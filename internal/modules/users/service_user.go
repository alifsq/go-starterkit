package user

import (
	"context"
	"errors"
	"fmt"

	"go-starterkit/internal/database/models"
	"go-starterkit/internal/response"

	"golang.org/x/crypto/bcrypt"
)

type ServiceUser interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error)
}

type serviceUser struct {
	repo RepositoryUser
}

func NewServiceUser(repo RepositoryUser) ServiceUser {
	return &serviceUser{
		repo: repo,
	}
}

func (s *serviceUser) CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	// Check duplicate email
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}
	if err != nil && !errors.Is(err, response.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Create user
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	if req.Address != "" {
		user.Address = &req.Address
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &CreateUserResponse{
		ID:    uint64(user.ID),
		Email: req.Email,
		Name:  req.Name,
		Role:  role,
	}, nil
}
