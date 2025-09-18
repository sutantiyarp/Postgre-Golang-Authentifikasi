package service

import (
	"hello-fiber/app/model"
	"hello-fiber/app/repository"
	"database/sql"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

// GetAllAlumniService untuk mengambil semua data alumni
func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
	alumniRepo := repository.NewAlumniRepository(db)
	alumniList, err := alumniRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    alumniList,
	})
}

// GetAlumniByIDService untuk mengambil data alumni berdasarkan ID
func GetAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	alumniRepo := repository.NewAlumniRepository(db)
	alumni, err := alumniRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    alumni,
	})
}

// CreateAlumniService untuk menambah alumni baru
func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "NIM, nama, jurusan, dan email harus diisi",
		})
	}

	alumniRepo := repository.NewAlumniRepository(db)
	alumni, err := alumniRepo.Create(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menambah alumni. Pastikan NIM dan email belum digunakan",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil ditambahkan",
		"data":    alumni,
	})
}

// UpdateAlumniService untuk mengupdate data alumni
func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Nama, jurusan, dan email harus diisi",
		})
	}

	alumniRepo := repository.NewAlumniRepository(db)
	alumni, err := alumniRepo.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil diupdate",
		"data":    alumni,
	})
}

// DeleteAlumniService untuk menghapus data alumni
func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	alumniRepo := repository.NewAlumniRepository(db)
	err = alumniRepo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus",
	})
}