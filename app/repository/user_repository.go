package repository

import (
	"hello-fiber/app/model"
	"database/sql"
	"log"
)

// UserRepository interface untuk operasional pengguna
type UserRepository interface {
	Save(user model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository membuat instance baru UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Save untuk menyimpan user baru
func (r *userRepository) Save(user model.User) (model.User, error) {
	sqlStatement := `
		INSERT INTO users (username, email, password, alumni_id, role_id, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int
	err := r.db.QueryRow(sqlStatement, user.Username, user.Email, user.Password, user.AlumniID, user.RoleID, user.CreatedAt).Scan(&id)
	if err != nil {
		log.Println("Error inserting user:", err)
		return model.User{}, err
	}
	user.ID = id
	return user, nil
}

// FindByEmail untuk mencari user berdasarkan email
func (r *userRepository) FindByEmail(email string) (model.User, error) {
	sqlStatement := `SELECT id, username, email, password, role_id, created_at FROM users WHERE email=$1`
	var user model.User
	err := r.db.QueryRow(sqlStatement, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt)
	if err != nil {
		log.Println("Error finding user by email:", err)
		return model.User{}, err
	}
	return user, nil
}