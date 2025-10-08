package service

import (
	"hello-fiber/app/model"
	"hello-fiber/app/repository"
	"database/sql"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
)

// GetAllPekerjaanAlumniWithPaginationService untuk mengambil data pekerjaan alumni dengan pagination, search, dan sorting
func GetAllPekerjaanAlumniWithPaginationService(c *fiber.Ctx, db *sql.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{
		"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "tanggal_mulai_kerja": true, "status_pekerjaan": true, "created_at": true, "is_delete": true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// Ambil data dari repository
	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaan, err := pekerjaanRepo.GetPekerjaanAlumniWithPagination(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pekerjaan alumni"})
	}

	total, err := pekerjaanRepo.CountPekerjaanAlumni(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count pekerjaan alumni"})
	}

	// Buat response pakai model
	response := model.PekerjaanAlumniResponse{
		Data: pekerjaan,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}

// GetAllPekerjaanAlumniService untuk mengambil semua data pekerjaan alumni
func GetAllPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	return GetAllPekerjaanAlumniWithPaginationService(c, db)
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

	if req.IsDelete == "" {
		req.IsDelete = "tidak"
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

// UpdatePekerjaanAlumniSementara untuk mengupdate data pekerjaan alumni
func UpdatePekerjaanAlumniSementara(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var req model.UpdatePekerjaanAlumniSoftDelete
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.IsDelete == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "is_delete harus diisi ('hapus' atau 'tidak')",
		})
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaan, err := pekerjaanRepo.Updatesementara(id, req)
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

	// Jika diset ke "hapus", tambahkan ke tabel trash (id parent = id pekerjaan_alumni)
	if strings.EqualFold(req.IsDelete, "hapus") {
		trashRepo := repository.NewPekerjaanAlumniTrashRepository(db)
		_ = trashRepo.InsertIfNotExists(id)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan alumni berhasil diupdate",
		"data":    pekerjaan,
	})
}

// UpdateDeleteStatusService untuk mengupdate status penghapusan pekerjaan alumni
func UpdateDeleteStatusService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var req model.UpdatePekerjaanAlumniSoftDelete
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	// Validasi input
	if req.IsDelete != "tidak" && req.IsDelete != "hapus" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Status is_delete harus 'tidak' atau 'hapus'",
		})
	}

	// Get user info from JWT token
	userID := c.Locals("user_id").(int)
	roleID := c.Locals("role_id").(int)

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	
	var pekerjaan model.PekerjaanAlumni
	
	// If admin, can update any record
	if roleID == 1 {
		pekerjaan, err = pekerjaanRepo.UpdateDeleteStatus(id, req.IsDelete)
	} else {
		// If user, can only update their own records
		pekerjaan, err = pekerjaanRepo.UpdateDeleteStatusByUser(id, userID, req.IsDelete)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Pekerjaan alumni tidak ditemukan atau Anda tidak memiliki akses",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate status penghapusan",
			"error":   err.Error(),
		})
	}

	// Tambahkan ke trash jika status 'hapus'
	if req.IsDelete == "hapus" {
		trashRepo := repository.NewPekerjaanAlumniTrashRepository(db)
		_ = trashRepo.InsertIfNotExists(pekerjaan.ID)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Status penghapusan berhasil diupdate",
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

// GetMyPekerjaanAlumniService untuk mengambil data pekerjaan alumni milik user yang sedang login
func GetMyPekerjaanAlumniService(c *fiber.Ctx, db *sql.DB) error {
	// Get user info from JWT token
	userID := c.Locals("user_id").(int)

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	pekerjaanList, err := pekerjaanRepo.GetByUserID(userID)
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

// GetPekerjaanAlumniTrashService untuk mengambil data trash
func GetPekerjaanAlumniTrashService(c *fiber.Ctx, db *sql.DB) error {
	trashRepo := repository.NewPekerjaanAlumniTrashRepository(db)
	list, err := trashRepo.List()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data trash", "error": err.Error()})
	}
	resp := model.PekerjaanAlumniTrashListResponse{Data: list}
	return c.JSON(resp)
}

// RestorePekerjaanAlumniTrashService untuk merestore data dari trash
func RestorePekerjaanAlumniTrashService(c *fiber.Ctx, db *sql.DB) error {
	trashID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID trash tidak valid"})
	}

	trashRepo := repository.NewPekerjaanAlumniTrashRepository(db)
	t, err := trashRepo.GetByTrashID(trashID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data trash tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data trash", "error": err.Error()})
	}

	pekerjaanRepo := repository.NewPekerjaanAlumniRepository(db)
	// Set is_delete menjadi 'tidak'
	_, err = pekerjaanRepo.UpdateDeleteStatus(t.PekerjaanAlumniID, "tidak")
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Pekerjaan alumni tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal merestore data", "error": err.Error()})
	}

	// Hapus dari trash
	if err := trashRepo.DeleteByTrashID(trashID); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus data dari trash", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil direstore dan dihapus dari trash"})
}

// HardDeletePekerjaanAlumniTrashService untuk menghapus data secara hard dari trash
func HardDeletePekerjaanAlumniTrashService(c *fiber.Ctx, db *sql.DB) error {
	trashID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID trash tidak valid"})
	}

	trashRepo := repository.NewPekerjaanAlumniTrashRepository(db)
	t, err := trashRepo.GetByTrashID(trashID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"success": false, "message": "Data trash tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data trash", "error": err.Error()})
	}

	// Gunakan transaksi agar konsisten
	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal memulai transaksi"})
	}
	commit := func() error { return tx.Commit() }
	rollback := func() { _ = tx.Rollback() }

	// Delete parent pekerjaan_alumni
	if _, err := tx.Exec(`DELETE FROM pekerjaan_alumni WHERE id = $1`, t.PekerjaanAlumniID); err != nil {
		rollback()
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal hard delete pekerjaan_alumni", "error": err.Error()})
	}

	// Delete dari trash
	if _, err := tx.Exec(`DELETE FROM pekerjaan_alumni_trash WHERE id = $1`, trashID); err != nil {
		rollback()
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus dari trash", "error": err.Error()})
	}

	if err := commit(); err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal commit transaksi", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data berhasil di-hard delete dan dihapus dari trash"})
}
