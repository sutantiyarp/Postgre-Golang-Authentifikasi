package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"password"`
	AlumniID  *bson.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id"`
	RoleID    bson.ObjectID `bson:"role_id" json:"role_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type Role struct {
	ID   bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Role string        `bson:"role" json:"role"`
}

type UserResponse struct {
	ID        bson.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string         `bson:"username" json:"username"`
	Email     string         `bson:"email" json:"email"`
	AlumniID  *bson.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id"`
	RoleID    bson.ObjectID  `bson:"role_id" json:"role_id"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
}

type RegisterRequest struct {
	Username  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"password"`
	RoleID    bson.ObjectID `bson:"role_id" json:"role_id"`
}

type CreateUserRequest struct {
	Username  string         `bson:"username" json:"username"`
	Email     string         `bson:"email" json:"email"`
	Password  string         `bson:"password" json:"password"`
	RoleID    bson.ObjectID  `bson:"role_id" json:"role_id"`
	AlumniID  *bson.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id"`
}

type UpdateUserRequest struct {
	Username  string         `bson:"username" json:"username"`
	Email     string         `bson:"email" json:"email"`
	Password  string         `bson:"password" json:"password"`
	RoleID    bson.ObjectID  `bson:"role_id" json:"role_id"`
	AlumniID  *bson.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id,omitempty"`
}

// LoginRequest untuk data input login
type LoginRequest struct {
	Email     string `bson:"username" json:"email"`
	Password  string `bson:"password" json:"password"`
}