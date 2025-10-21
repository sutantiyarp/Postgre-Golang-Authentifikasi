package service

import (
	"hello-fiber/app/model"
	"hello-fiber/app/repository"
	"hello-fiber/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

var userRepo = repository.NewUserRepositoryMongo()

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var req model.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username, email, dan password harus diisi"})
	}

	if !isValidEmail(req.Email) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format email tidak valid"})
	}

	alumniID, err := bson.ObjectIDFromHex(req.AlumniID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Alumni ID tidak valid"})
	}

	roleID, err := bson.ObjectIDFromHex(req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Role ID tidak valid"})
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		AlumniID: alumniID,
		RoleID:   roleID,
	}

	id, err := userRepo.CreateUser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mendaftarkan user", "error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "User berhasil didaftarkan", "id": id.Hex()})
}

func LoginService(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email dan password harus diisi"})
	}

	user, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Email atau password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Email atau password salah"})
	}

	token, err := utils.GenerateJWT(*user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal membuat token", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Login berhasil", "token": token, "user": user})
}
