package user

import (
	"context"
	"errors"
	"fmt"

	"go-starterkit/internal/database/models"
	"go-starterkit/internal/modules/users/dto"
	"go-starterkit/internal/modules/users/mapper"
	"go-starterkit/internal/response"

	"golang.org/x/crypto/bcrypt"
)

type ServiceUser interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetMyProfile(ctx context.Context, id int64) (*models.User, error)
	GetAllUsersPaginate(ctx context.Context, page, pageSize int) ([]dto.UserResponse, int64, error)
}

type serviceUser struct {
	repo RepositoryUser
}

func NewServiceUser(repo RepositoryUser) ServiceUser {
	return &serviceUser{
		repo: repo,
	}
}

func (s *serviceUser) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
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

	return &dto.CreateUserResponse{
		ID:    uint64(user.ID),
		Email: req.Email,
		Name:  req.Name,
		Role:  role,
	}, nil
}

func (s *serviceUser) GetMyProfile(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *serviceUser) GetAllUsersPaginate(ctx context.Context, page, pageSize int) ([]dto.UserResponse, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default 10 data per halaman
	}

	users, total, err := s.repo.GetUsersPaginate(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Konversi ke DTO menggunakan helper yang sudah kita buat sebelumnya
	response := mapper.ToUserResponseList(users)
	return response, total, nil
}
