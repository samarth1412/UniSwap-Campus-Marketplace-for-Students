package models

type RegisterRequest struct {
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	University string `json:"university"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
