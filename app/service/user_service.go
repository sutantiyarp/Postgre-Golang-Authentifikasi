package service

import (
	"hello-fiber/app/model"
	"hello-fiber/app/repository"
	"hello-fiber/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"unicode"
)

var userRepo = repository.NewUserRepositoryMongo()

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	for _, char := range username {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
	}
	return true
}

func isValidPassword(password string) bool {
	if len(password) < 5 {
		return false
	}

	hasUpper := true
	hasLower := true
	hasNumber := true

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
	}

	return hasUpper && hasLower && hasNumber
}

func toUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AlumniID:  user.AlumniID,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
	}
}

func isObjectIDEmpty(id *bson.ObjectID) bool {
	return id == nil || id.IsZero()
}

func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username, email, dan password harus diisi"})
	}

	if !isValidUsername(req.Username) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username harus 3-50 karakter, hanya alphanumeric dan underscore"})
	}

	if !isValidEmail(req.Email) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format email tidak valid"})
	}

	if !isValidPassword(req.Password) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Password minimal 5 karakter"})
	}

	existingUser, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal validasi username", "error": err.Error()})
	}
	if existingUser != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username sudah terdaftar"})
	}

	if req.RoleID == bson.NilObjectID {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Role ID harus diisi"})
	}

	id, err := userRepo.Register(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mendaftarkan user", "error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "User berhasil didaftarkan", "id": id.Hex()})
}

func CreateUserAdmin(c *fiber.Ctx) error {
	var req model.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username, email, dan password harus diisi"})
	}

	if !isValidUsername(req.Username) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username harus 3-50 karakter, hanya alphanumeric dan underscore"})
	}

	if !isValidEmail(req.Email) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format email tidak valid"})
	}

	if !isValidPassword(req.Password) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Password minimal 8 karakter, harus ada uppercase, lowercase, dan number"})
	}

	existingUser, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal validasi username", "error": err.Error()})
	}
	if existingUser != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username sudah terdaftar"})
	}

	if req.RoleID == bson.NilObjectID {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Role ID harus diisi"})
	}

	_, err = userRepo.GetRoleByID(req.RoleID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Role ID tidak valid", "error": err.Error()})
	}

	id, err := userRepo.CreateUser(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal membuat user", "error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "User berhasil dibuat", "id": id.Hex()})
}

func LoginService(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Email dan password harus diisi"})
	}

	user, err := userRepo.GetUserByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
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

	return c.JSON(fiber.Map{"success": true, "message": "Login berhasil", "token": token, "user": toUserResponse(user)})
}

func GetAllUsersService(c *fiber.Ctx) error {
	page := int64(1)
	limit := int64(10)

	if p := c.Query("page"); p != "" {
		page = int64(c.QueryInt("page", 1))
	}
	if l := c.Query("limit"); l != "" {
		limit = int64(c.QueryInt("limit", 10))
	}

	users, total, err := userRepo.GetAllUsers(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data user", "error": err.Error()})
	}

	// Convert ke UserResponse untuk hide password
	var userResponses []model.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *toUserResponse(&user))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data user berhasil diambil",
		"data":    userResponses,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func UpdateUserService(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req model.UpdateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	hasUpdate := req.Username != "" || req.Email != "" || req.Password != "" || req.RoleID != bson.NilObjectID || req.AlumniID != nil
	if !hasUpdate {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Minimal ada satu field yang harus diupdate (username, email, password, role_id, atau alumni_id)"})
	}

	// Validasi input jika ada
	if req.Username != "" && !isValidUsername(req.Username) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username harus 3-50 karakter, hanya alphanumeric dan underscore"})
	}

	if req.Email != "" && !isValidEmail(req.Email) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format email tidak valid"})
	}

	if req.Password != "" && !isValidPassword(req.Password) {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Password minimal 5 karakter"})
	}

	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "User ID tidak valid"})
	}

	// Cek apakah username sudah ada (jika diupdate)
	if req.Username != "" {
		existingUser, err := userRepo.GetUserByUsername(req.Username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal validasi username", "error": err.Error()})
		}
		if existingUser != nil && existingUser.ID != id {
			return c.Status(400).JSON(fiber.Map{"success": false, "message": "Username sudah terdaftar"})
		}
	}

	if err := userRepo.UpdateUser(id, req); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal update user", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User berhasil diupdate"})
}

func DeleteUserService(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "User ID harus diisi"})
	}

	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "User ID tidak valid"})
	}

	if err := userRepo.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal delete user", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "User berhasil dihapus"})
}