package response

import "errors"

var(
	ErrForbidden = errors.New("forbidden")
	ErrEventNotFound = errors.New("event not found")
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExist = errors.New("email already registered")
)