package auth

// SignupRequest represents user registration input
type SignupRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// LoginRequest represents user login input
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
