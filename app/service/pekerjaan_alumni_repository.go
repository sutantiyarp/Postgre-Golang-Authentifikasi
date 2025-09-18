package service

import (
	"hello-fiber/app/model"
	"hello-fiber/app/repository"
	"database/sql"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

// GetAllPekerjaanAlumniService untuk mengambil semua data pekerjaan alumni
func GetAllPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaanList, err := pekerjaanRepo.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diambil",
		"data":    pekerjaanList,
	})
}

// GetPekerjaanAlumniByIDService untuk mengambil data pekerjaan alumni berdasarkan ID
func GetPekerjaanAlumniByIDService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaan, err := pekerjaanRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Pekerjaan alumni tidak ditemukan",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diambil",
		"data":    pekerjaan,
	})
}

// GetPekerjaanAlumniByAlumniIDService untuk mengambil semua pekerjaan berdasarkan alumni ID
func GetPekerjaanAlumniByAlumniIDService(c *fiber.Ctx, db *sql.DB) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Alumni ID tidak valid",
		})
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaanList, err := pekerjaanRepo.GetByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan alumni berhasil diambil",
		"data":    pekerjaanList,
	})
}

// CreatePekerjaanAlumniService untuk menambah pekerjaan alumni baru
func CreatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.AlumniID == 0 || req.NamaPerusahaan == "" || req.PosisiJabatan == "" || 
	   req.BidangIndustri == "" || req.LokasiKerja == "" || req.TanggalMulaiKerja == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Alumni ID, nama perusahaan, posisi jabatan, bidang industri, lokasi kerja, dan tanggal mulai kerja harus diisi",
		})
	}

	// Set default status jika kosong
	if req.StatusPekerjaan == "" {
		req.StatusPekerjaan = "aktif"
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaan, err := pekerjaanRepo.Create(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menambah pekerjaan alumni",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan alumni berhasil ditambahkan",
		"data":    pekerjaan,
	})
}

// UpdatePekerjaanAlumniService untuk mengupdate data pekerjaan alumni
func UpdatePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var req model.UpdatePekerjaanAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.NamaPerusahaan == "" || req.PosisiJabatan == "" || 
	   req.BidangIndustri == "" || req.LokasiKerja == "" || req.TanggalMulaiKerja == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Nama perusahaan, posisi jabatan, bidang industri, lokasi kerja, dan tanggal mulai kerja harus diisi",
		})
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaan, err := pekerjaanRepo.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Pekerjaan alumni tidak ditemukan",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate pekerjaan alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan alumni berhasil diupdate",
		"data":    pekerjaan,
	})
}

// DeletePekerjaanAlumniService untuk menghapus data pekerjaan alumni
func DeletePekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	err = pekerjaanRepo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Pekerjaan alumni tidak ditemukan",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus pekerjaan alumni",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan alumni berhasil dihapus",
	})
}