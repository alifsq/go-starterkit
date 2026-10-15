package user

type CreateUserResponse struct {
    ID    uint64 `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
    Role  string `json:"role"`
}