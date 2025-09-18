package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	RoleID    int       `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
}

// UserRequest untuk data input registrasi
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

// LoginRequest untuk data input login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}