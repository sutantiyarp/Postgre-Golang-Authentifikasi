package service

import (
	"strconv"
	"hello-fiber/app/model"
	"hello-fiber/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var alumniRepo = repository.NewAlumniRepositoryMongo()

// GetAllAlumniService mengambil semua alumni
func GetAllAlumniService(c *fiber.Ctx) error {
	// Parse query params
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	alumni, total, err := alumniRepo.GetAllAlumniWithPagination(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data alumni", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data alumni",
		"data":    alumni,
		"page":    page,
		"limit":   limit,
		"total":   total,
	})
}

// GetAlumniByIDService mengambil alumni berdasarkan id (bson.ObjectID)
func GetAlumniByIDService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format ID tidak valid"})
	}

	item, err := alumniRepo.GetAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data alumni", "data": item})
}

// CreateAlumniService menambah alumni baru
func CreateAlumniService(c *fiber.Ctx) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	// validasi minimal
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIM, nama, jurusan, dan email harus diisi"})
	}

	id, err := alumniRepo.CreateAlumni(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menambahkan alumni", "error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Alumni berhasil dibuat", "id": id.Hex()})
}

// UpdateAlumniService mengupdate alumni berdasarkan id
func UpdateAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format ID tidak valid"})
	}

	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	hasUpdate := req.NIM != nil || req.Nama != nil || req.Jurusan != nil || req.Angkatan != nil || 
		req.TahunLulus != nil || req.Email != nil || req.NoTelepon != nil || req.Alamat != nil
	
	if !hasUpdate {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Minimal ada satu field yang harus diupdate"})
	}

	if err := alumniRepo.UpdateAlumni(id, req); err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengupdate alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil diupdate"})
}

// DeleteAlumniService menghapus alumni berdasarkan id
func DeleteAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Format ID tidak valid"})
	}

	if err := alumniRepo.DeleteAlumni(id); err != nil {
		// Check if it's a "not found" error
		if err.Error() == "alumni not found" {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan", "error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil dihapus"})
}