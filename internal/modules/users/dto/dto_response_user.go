package dto

import "time"

type CreateUserResponse struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type MyProfileResponse struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	IsEmailVerified bool       `json:"is_email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Phone           *string    `json:"phone"`
	Address         *string    `json:"address"`
	Role            string     `json:"role"`
	CreatedAt       time.Time  `json:"created_at"`
}
type UserResponse struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	IsEmailVerified bool       `json:"is_email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	Phone           *string    `json:"phone"`
	Address         *string    `json:"address"`
	Role            string     `json:"role"`
	CreatedAt       time.Time  `json:"created_at"`
}