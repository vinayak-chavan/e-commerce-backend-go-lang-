package dto

// UserRegisterRequest represents the request body for user registration
type UserRegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    Role     string `json:"role" binding:"required"`
}

// UserRegisterResponse represents the response body for user registration
type UserRegisterResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

// UserLoginRequest represents the request body for user login
type UserLoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// UserLoginResponse represents the response body for user login
type UserLoginResponse struct {
    User  UserRegisterResponse `json:"user"`
    Token string               `json:"token"`
}

type SuccessResponse struct {
    Message string `json:"message"`
}