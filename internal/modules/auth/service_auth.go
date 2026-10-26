package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"go-starterkit/internal/database/models"
	user "go-starterkit/internal/modules/users"
	"go-starterkit/internal/package/jwt"
	"go-starterkit/internal/package/logger"
	"go-starterkit/internal/package/mail"
	"go-starterkit/internal/package/utils"
	"go-starterkit/internal/response"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ServiceAuth interface {
	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, email, password, userAgent, ipAddress string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	LogoutAllDevices(ctx context.Context, userID int64) error
	GetAuthenticatedUser(ctx context.Context, userID int64) (*models.User, error)
	CleanupExpiredTokens(ctx context.Context) error

	// Email
	VerifyEmail(ctx context.Context, plainToken string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, plainToken, newPassword string) error
	ResendVerificationEmail(ctx context.Context, userID int64) error
}

type serviceAuth struct {
	authRepo RepositoryAuth
	userRepo user.RepositoryUser
	jwtSvc   jwt.JWTService
	mail     mail.Mailer
	log      logger.Logger
}

func NewServiceAuth(authRepo RepositoryAuth, userRepo user.RepositoryUser, jwtSvc jwt.JWTService, mail mail.Mailer, log logger.Logger) ServiceAuth {
	return &serviceAuth{authRepo: authRepo, userRepo: userRepo, jwtSvc: jwtSvc, mail: mail, log: log}
}

func generateToken() (plainToken string, tokenHash string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", nil
	}
	plainToken = hex.EncodeToString(bytes)
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash = hex.EncodeToString(hash[:])
	return plainToken, tokenHash, nil
}

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must contain at least 8 characters")
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSymbol {
		return errors.New("password must include uppercase, lowercase, number, and symbol")
	}
	return nil
}

func (s *serviceAuth) Register(ctx context.Context, req *RegisterRequest) error {
	if req == nil {
		return errors.New("request is required")
	}

	if err := validatePasswordStrength(req.Password); err != nil {
		s.log.Warn("Password strength validation failed", zap.String("email", req.Email), zap.Error(err))
		return err
	}
	normalizedEmail := strings.ToLower(req.Email)
	existingUser, err := s.userRepo.FindByEmail(ctx, normalizedEmail)
	if err == nil && existingUser != nil {
		s.log.Warn("Registration Failed : email already exists")
		return response.ErrEmailExist
	}
	if err != nil && !errors.Is(err, response.ErrUserNotFound) {
		s.log.Error("Failed to check existing user in database", zap.String("email", normalizedEmail), zap.Error(err))
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed To hash password", zap.String("email", normalizedEmail), zap.Error(err))
		return fmt.Errorf("failed to hash password: %w", err)
	}

	userModel := &models.User{
		Name:            req.Name,
		Email:           strings.ToLower(req.Email),
		PasswordHash:    string(hash),
		Role:            "user",
		IsActive:        true,
		IsEmailVerified: false,
	}

	// create user
	if _, err := s.authRepo.Register(ctx, userModel); err != nil {
		s.log.Error("Failed To Create User", zap.Error(err))
		return err
	}

	// Generate Token Verifikasi Email
	plainToken, tokenHash, err := generateToken()
	if err != nil {
		s.log.Error("Failed to generate verification token", zap.Error(err))
		return err
	}

	tokenModel := &models.EmailVerificationToken{
		UserID:    userModel.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Expire 24 jam
	}

	if err := s.authRepo.SaveEmailVerificationToken(ctx, tokenModel); err != nil {
		s.log.Error("Failed to save verification token", zap.Error(err))
		return err
	}

	userEmail := userModel.Email
	go func() {
		if err := s.mail.SendVerificationEmail(context.Background(), userEmail, plainToken); err != nil {
			s.log.Error("Background task: Failed to send verification email", zap.Error(err), zap.String("email", userEmail))
		}
	}()

	return nil
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

// Email
func (s *serviceAuth) VerifyEmail(ctx context.Context, plainToken string) error {
	if plainToken == "" {
		return errors.New("Token Verification Is Required")
	}

	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	tokenModel, err := s.authRepo.FindEmailVerificationToken(ctx, tokenHash)
	if err != nil || tokenModel == nil {
		s.log.Warn("Invalid email Verification token attempt", zap.String("token_hash", tokenHash))
		return errors.New("Invalid Or Expired Verification Token")
	}

	// validasi status token
	if tokenModel.UsedAt != nil {
		return errors.New("token has already been used")
	}
	if time.Now().After(tokenModel.ExpiresAt) {
		return errors.New("Token Has Expired")
	}

	// mark token `used`
	if err := s.authRepo.MarkEmailVerificationTokenAsUsed(ctx, tokenModel.ID); err != nil {
		s.log.Error("Failed to mark verification token as used", zap.Error(err))
		return err
	}
	// Update kolom is_email_verified = true di tabel users
	if err := s.authRepo.MarkEmailVerificationTokenAsUsed(ctx, tokenModel.UserID); err != nil {
		s.log.Error("Failed to update user verification status", zap.Error(err), zap.Int64("user_id", tokenModel.UserID))
		return err
	}

	if err := s.authRepo.MarkUserEmailAsVerified(ctx, tokenModel.UserID); err != nil {
		s.log.Error("Failed to Update Verified user status", zap.Error(err), zap.Int64("user_id", tokenModel.UserID))
		return err
	}

	s.log.Info("Email verified successfully", zap.Int64("user_id", tokenModel.UserID))
	return nil
}
func (s *serviceAuth) ForgotPassword(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		// Security Best Practice: Tetap return nil agar attacker tidak bisa menebak email mana yang registered
		s.log.Info("Forgot password requested for non-existing email", zap.String("email", email))
		return nil
	}

	// Generate reset token
	plainToken, tokenHash, err := generateToken()
	if err != nil {
		s.log.Error("Failed to generate reset password token", zap.Error(err))
		return err
	}

	tokenModel := &models.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * time.Minute), // Expire pendek (15 Menit)
	}

	if err := s.authRepo.SavePasswordResetToken(ctx, tokenModel); err != nil {
		s.log.Error("Failed to save password reset token", zap.Error(err))
		return err
	}

	// Kirim Email Reset Password secara Async (Goroutine)
	userEmail := user.Email
	go func() {
		if err := s.mail.SendResetPasswordEmail(context.Background(), userEmail, plainToken); err != nil {
			s.log.Error("Failed to send reset password email", zap.Error(err), zap.String("email", userEmail))
		}
	}()

	return nil
}

func (s *serviceAuth) ResetPassword(ctx context.Context, plainToken, newPassword string) error {
	if plainToken == "" || newPassword == "" {
		return errors.New("token and new password are required")
	}

	// Hash token request
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Cari token di DB
	tokenModel, err := s.authRepo.FindPasswordResetToken(ctx, tokenHash)
	if err != nil || tokenModel == nil {
		return errors.New("invalid or expired reset token")
	}

	// Validasi expired & state token
	if tokenModel.UsedAt != nil {
		return errors.New("token has already been used")
	}
	if time.Now().After(tokenModel.ExpiresAt) {
		return errors.New("token has expired")
	}

	// validate Strenght Password
	if err := validatePasswordStrength(newPassword); err != nil {
		s.log.Warn("Password strength validation failed", zap.Error(err))
		return err
	}

	// Hash password baru (contoh menggunakan bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password di tabel users & mark token used
	if err := s.userRepo.UpdateUserPassword(ctx, tokenModel.UserID, string(hashedPassword)); err != nil {
		return err
	}

	if err := s.authRepo.MarkPasswordResetTokenAsUsed(ctx, tokenModel.ID); err != nil {
		s.log.Error("Failed to mark reset token as used", zap.Error(err))
		return err
	}

	s.log.Info("Password reset successful", zap.Int64("user_id", tokenModel.UserID))
	return nil
}

func (s *serviceAuth) ResendVerificationEmail(ctx context.Context, userID int64) error {
	// Find User ID check verification user
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("User NotFound")
	}

	if user.IsEmailVerified {
		return errors.New("Account verified")
	}

	// GenerateNewToken
	plainToken, tokenHash, err := generateToken()
	if err != nil {
		s.log.Error("Failed to generate verification token", zap.Error(err))
		return err
	}

	tokenModel := &models.EmailVerificationToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Expire 24 jam
	}

	// update token verifikasi
	err = s.authRepo.UpdateEmailVerificationToken(ctx, tokenModel)
	if err != nil {
		s.log.Error("Failed Update Token Verification Email", zap.Error(err))
		return err
	}

	// send email
	go func() {
		if err := s.mail.SendVerificationEmail(ctx, user.Email, plainToken); err != nil {
			s.log.Error("Background task: Failed to send verification email", zap.Error(err), zap.String("email", user.Email))
		}
	}()
	return nil
}