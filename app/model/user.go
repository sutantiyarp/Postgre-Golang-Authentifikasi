package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type User struct {
	ID        bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string         `bson:"username" json:"username"`
	Email     string         `bson:"email" json:"email"`
	Password  string         `bson:"password" json:"password"`
	AlumniID  *bson.ObjectID `bson:"alumni_id" json:"alumni_id,omitempty"` // Menggunakan pointer untuk bisa menjadi nil
	RoleID    bson.ObjectID  `bson:"role_id" json:"role_id"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
}

type Role struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Role string        `bson:"role" json:"role"`
}

// UserRequest untuk data input registrasi
type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AlumniID  string `json:"alumni_id,omitempty"`
	RoleID    string `json:"role_id"`
}

// LoginRequest untuk data input login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
