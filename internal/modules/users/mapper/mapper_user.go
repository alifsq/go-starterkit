package mapper

import (
	"go-starterkit/internal/database/models"
	"go-starterkit/internal/modules/users/dto"
)

func ToMyProfileResponse(user *models.User) *dto.MyProfileResponse {
	return &dto.MyProfileResponse{
		ID:              uint(user.ID),
		Name:            user.Name,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Phone:           user.Phone,
		Address:         user.Address,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt,
	}
}

func ToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              uint(user.ID),
		Name:            user.Name,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Phone:           user.Phone,
		Address:         user.Address,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt,
	}
}

func ToUserResponseList(users []*models.User) []dto.UserResponse {
	var list []dto.UserResponse
	for _, user := range users {
		list = append(list, *ToUserResponse(user))
	}
	return list
}
