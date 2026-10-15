package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go-starterkit/internal/database/models"
	user "go-starterkit/internal/modules/users"
	"go-starterkit/internal/package/jwt"
	"go-starterkit/internal/package/utils"
	"go-starterkit/internal/response"

	"golang.org/x/crypto/bcrypt"
)

type ServiceAuth interface {
	Register(ctx context.Context, req *RegisterRequest) (*models.User, error)
	Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	LogoutAllDevices(ctx context.Context, userID int64) error
	GetAuthenticatedUser(ctx context.Context, userID int64) (*models.User, error)
	CleanupExpiredTokens(ctx context.Context) error
}

type serviceAuth struct {
	authRepo RepositoryAuth
	userRepo user.RepositoryUser
	jwtSvc   jwt.JWTService
}

func NewServiceAuth(authRepo RepositoryAuth, userRepo user.RepositoryUser, jwtSvc jwt.JWTService) ServiceAuth {
	return &serviceAuth{authRepo: authRepo, userRepo: userRepo, jwtSvc: jwtSvc}
}

func (s *serviceAuth) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, strings.ToLower(req.Email))
	if err == nil && existingUser != nil {
		return nil, response.ErrEmailExist
	}
	if err != nil && !errors.Is(err, response.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userModel := &models.User{
		Name:         req.Name,
		Email:        strings.ToLower(req.Email),
		PasswordHash: string(hash),
		Role:         "user",
		IsActive:     true,
	}
	if _, err := s.authRepo.Register(ctx, userModel); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return userModel, nil
}

func tokenTTLFromEnv(key string, fallback time.Duration) time.Duration {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func (s *serviceAuth) Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return "", "", errors.New("Invalid email or password")
	}

	userModel, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("Invalid email or password")
	}
	if !userModel.IsActive {
		return "", "", errors.New("Invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("Invalid email or password")
	}

	accessTTL := tokenTTLFromEnv("JWT_ACCESS_TOKEN_EXP", 15*time.Minute)
	refreshTTL := tokenTTLFromEnv("JWT_REFRESH_TOKEN_EXP", 24*time.Hour)

	accessToken, err := s.jwtSvc.GenerateAccessToken(userModel.ID, userModel.Role, accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, _, err := s.jwtSvc.GenerateRefreshToken(userModel.ID, refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	userAgentValue := userAgent
	ipAddressValue := ipAddress
	tokenRecord := &models.RefreshToken{
		UserID:    userModel.ID,
		TokenHash: utils.HashToken(refreshToken),
		UserAgent: &userAgentValue,
		IPAddress: &ipAddressValue,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(refreshTTL),
	}
	if _, err := s.authRepo.SaveRefreshToken(ctx, tokenRecord); err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
			_ = s.authRepo.RevokeToken(ctx, utils.HashToken(refreshToken))
		}
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

func (s *serviceAuth) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}
	return s.authRepo.RevokeToken(ctx, utils.HashToken(refreshToken))
}

func (s *serviceAuth) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return "", "", errors.New("refresh token not found")
	}

	claims, err := s.jwtSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	storedToken, err := s.authRepo.FindRefreshTokenByHash(ctx, utils.HashToken(refreshToken))
	if err != nil || storedToken == nil {
		return "", "", errors.New("refresh token is invalid or revoked")
	}
	if storedToken.IsRevoked || storedToken.ExpiresAt.Before(time.Now()) {
		return "", "", errors.New("refresh token is invalid or revoked")
	}

	userModel, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil || userModel == nil || !userModel.IsActive {
		return "", "", errors.New("user is not active")
	}

	accessTTL := tokenTTLFromEnv("JWT_ACCESS_TOKEN_EXP", 15*time.Minute)
	refreshTTL := tokenTTLFromEnv("JWT_REFRESH_TOKEN_EXP", 24*time.Hour)

	newAccessToken, err := s.jwtSvc.GenerateAccessToken(userModel.ID, userModel.Role, accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	newRefreshToken, _, err := s.jwtSvc.GenerateRefreshToken(userModel.ID, refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.authRepo.RevokeToken(ctx, utils.HashToken(refreshToken)); err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}
	if _, err := s.authRepo.SaveRefreshToken(ctx, &models.RefreshToken{
		UserID:    userModel.ID,
		TokenHash: utils.HashToken(newRefreshToken),
		UserAgent: storedToken.UserAgent,
		IPAddress: storedToken.IPAddress,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(refreshTTL),
	}); err != nil {
		return "", "", fmt.Errorf("failed to store rotated refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *serviceAuth) LogoutAllDevices(ctx context.Context, userID int64) error {
	return s.authRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *serviceAuth) GetAuthenticatedUser(ctx context.Context, userID int64) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *serviceAuth) CleanupExpiredTokens(ctx context.Context) error {
	return s.authRepo.DeleteExpiredTokens(ctx)
}
