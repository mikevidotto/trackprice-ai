package auth

// SignupRequest represents user registration input
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents user login input
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
