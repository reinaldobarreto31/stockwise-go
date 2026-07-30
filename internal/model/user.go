package model

import "time"

// User represents a system user.
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // admin | operator
	CreatedAt    time.Time `json:"created_at"`
}

// RegisterInput is the request body for user registration.
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput is the request body for login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is the response for a successful login.
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
