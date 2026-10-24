package dto


type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=100"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Phone    string `json:"phone" binding:"omitempty"`
    Address  string `json:"address" binding:"omitempty"`
    Role     string `json:"role" binding:"omitempty,oneof=admin organizer user moderator"`
}

