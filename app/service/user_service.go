package service

import (
	"hello-fiber/app/model"
	"hello-fiber/app/repository"
	"log"
	"time"
	"golang.org/x/crypto/bcrypt"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	// "hello-fiber/utils"
)

func RegisterService(c *fiber.Ctx, db *sql.DB) error {
	// Parse request body
	var userRequest model.UserRequest
	if err := c.BodyParser(&userRequest); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(500).SendString("Error processing password")
	}

	// Siapkan user model
	user := model.User{
		Username:  userRequest.Username,
		Email:     userRequest.Email,
		Password:  string(hashedPassword),
		RoleID:    userRequest.RoleID,
		CreatedAt: time.Now(),
	}

	// Simpan user ke database
	userRepo := repository.NewUserRepository(db)
	createdUser, err := userRepo.Save(user)
	if err != nil {
		log.Println("Error saving user:", err)
		return c.Status(500).SendString("Error registering user")
	}

	// Response sukses
	return c.Status(201).JSON(createdUser)
}

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	// Parse request body untuk login
	var loginData model.LoginRequest
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	// Membuat repository user
	userRepo := repository.NewUserRepository(db)

	// Cek email dan password
	user, err := userRepo.FindByEmail(loginData.Email)
	if err != nil {
		return c.Status(401).SendString("Invalid credentials")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		return c.Status(401).SendString("Invalid credentials")
	}

	// Membuat klaim token JWT
	claims := jwt.MapClaims{
		"sub": user.ID,  // ID pengguna
		"exp": time.Now().Add(time.Hour * 24).Unix(),  // Waktu kadaluarsa token (1 hari)
		"iat": time.Now().Unix(),  // Waktu pembuatan token
		"role": user.RoleID,  // Peran pengguna
	}

	// Membuat token dengan klaim dan menandatangani menggunakan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_jwt_secret"))
	if err != nil {
		return c.Status(500).SendString("Error generating token")
	}

	// Mengembalikan token ke client
	return c.Status(200).JSON(fiber.Map{"token": tokenString})
}