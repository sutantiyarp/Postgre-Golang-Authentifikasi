package service

import (
	"errors"

	"hello-fiber/app/model"
	"hello-fiber/app/repository"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var pekerjaanRepo = repository.NewPekerjaanAlumniRepositoryMongo()

// GetAllPekerjaanAlumniService mengambil semua pekerjaan alumni (filter sederhana dengan query is_delete/search opsional)
func GetAllPekerjaanAlumniService(c *fiber.Ctx) error {
	// optional search query handled by repo if implemented; here we return all non-deleted
	data, err := pekerjaanRepo.GetAllPekerjaanAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data pekerjaan alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data pekerjaan alumni", "data": data})
}

// GetPekerjaanAlumniByIDService mengambil pekerjaan berdasarkan id (bson.ObjectID)
func GetPekerjaanAlumniByIDService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	item, err := pekerjaanRepo.GetPekerjaanAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Pekerjaan alumni tidak ditemukan", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data pekerjaan alumni", "data": item})
}

// GetPekerjaanAlumniByAlumniIDService mengambil pekerjaan berdasarkan alumni_id (bson.ObjectID)
func GetPekerjaanAlumniByAlumniIDService(c *fiber.Ctx) error {
	alumniIDStr := c.Params("alumni_id")
	if alumniIDStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "alumni_id tidak boleh kosong"})
	}

	alumniID, err := bson.ObjectIDFromHex(alumniIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "alumni_id format tidak valid"})
	}

	list, err := pekerjaanRepo.GetPekerjaanAlumniByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data pekerjaan alumni dari alumni_id", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data pekerjaan alumni dari alumni_id", "data": list})
}

// CreatePekerjaanAlumniService menambah pekerjaan alumni baru
func CreatePekerjaanAlumniService(c *fiber.Ctx) error {
    var req model.CreatePekerjaanAlumniRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
    }

    if err := validatePekerjaanInput(req); err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
    }

	if req.StatusPekerjaan == "" {
		req.StatusPekerjaan = "aktif"
	}

	id, err := pekerjaanRepo.CreatePekerjaanAlumni(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menambahkan pekerjaan alumni", "error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Pekerjaan alumni berhasil dibuat", "id": id})
}

// UpdatePekerjaanAlumniService mengupdate pekerjaan alumni berdasarkan id
func UpdatePekerjaanAlumniService(c *fiber.Ctx) error {
    idStr := c.Params("id")
    if idStr == "" {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
    }

    id, err := bson.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
    }

    var req model.UpdatePekerjaanAlumniRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
    }

    if err := validatePekerjaanUpdateInput(req); err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
    }

    err = pekerjaanRepo.UpdatePekerjaanAlumni(id, req)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengupdate pekerjaan alumni", "error": err.Error()})
    }

    return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan alumni berhasil diupdate"})
}


// SoftDeletePekerjaanAlumniService set is_delete pada dokumen (soft delete)
func SoftDeletePekerjaanAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	var body struct {
		IsDelete string `json:"is_delete"`
	}
	
	// Parse body, default to "hapus" (soft delete)
	if err := c.BodyParser(&body); err != nil {
		body.IsDelete = "hapus"
	}

	if body.IsDelete != "hapus" && body.IsDelete != "tidak" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "is_delete harus bernilai 'hapus' atau 'tidak'"})
	}

	isDelete := body.IsDelete == "hapus"
	if err := pekerjaanRepo.SoftDeletePekerjaanAlumni(id, isDelete); err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengubah status is_delete", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Status is_delete berhasil diupdate"})
}

// GetTrashedPekerjaanAlumniService mengambil semua yang is_delete == true
func GetTrashedPekerjaanAlumniService(c *fiber.Ctx) error {
	data, err := pekerjaanRepo.GetTrashedPekerjaanAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil trashed items", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil trashed items", "data": data})
}

// HardDeleteTrashedPekerjaanAlumniService menghapus permanen item yang status hapus dan mengembalikan dokumen yang dihapus
func HardDeleteTrashedPekerjaanAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	item, err := pekerjaanRepo.HardDeleteTrashedPekerjaanAlumni(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal menghapus permanen atau item tidak ditemukan", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil menghapus permanen", "data": item})
}

// RestoreTrashedPekerjaanAlumniService me-restore item dari trash
func RestoreTrashedPekerjaanAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	item, err := pekerjaanRepo.RestoreTrashedPekerjaanAlumni(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal merestore item atau item tidak ditemukan", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Berhasil merestore item", "data": item})
}

// DeletePekerjaanAlumniService - wrapper for soft delete with is_delete = true
func DeletePekerjaanAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	if err := pekerjaanRepo.SoftDeletePekerjaanAlumni(id, true); err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal menghapus pekerjaan alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan alumni berhasil dihapus"})
}

func validatePekerjaanInput(p model.CreatePekerjaanAlumniRequest) error {
	if p.AlumniID == bson.NilObjectID {
		return errors.New("alumni_id wajib diisi dan harus valid")
	}
	if p.NamaPerusahaan == "" {
		return errors.New("nama_perusahaan wajib diisi")
	}
	if p.PosisiJabatan == "" {
		return errors.New("posisi_jabatan wajib diisi")
	}
	if p.LokasiKerja == "" {
		return errors.New("lokasi_kerja wajib diisi")
	}
	if p.TanggalMulaiKerja.IsZero() {
		return errors.New("tanggal_mulai_kerja wajib diisi")
	}
	
	if !p.TanggalSelesaiKerja.IsZero() && p.TanggalSelesaiKerja.Before(p.TanggalMulaiKerja.Time) {
		return errors.New("tanggal_selesai_kerja tidak boleh lebih awal dari tanggal_mulai_kerja")
	}
	
	return nil
}

// service/pekerjaan_alumni_service.go

func validatePekerjaanUpdateInput(u model.UpdatePekerjaanAlumniRequest) error {
    nothingToUpdate :=
        u.AlumniID == nil &&
        u.NamaPerusahaan == "" &&
        u.PosisiJabatan == "" &&
        u.BidangIndustri == "" &&
        u.LokasiKerja == "" &&
        u.GajiRange == "" &&
        u.StatusPekerjaan == "" &&
        u.DeskripsiPekerjaan == "" &&
        u.TanggalMulaiKerja == nil &&
        u.TanggalSelesaiKerja == nil

    if nothingToUpdate {
        return errors.New("tidak ada field yang diupdate")
    }
    
    if u.TanggalMulaiKerja != nil && u.TanggalSelesaiKerja != nil &&
        u.TanggalSelesaiKerja.Time.Before(u.TanggalMulaiKerja.Time) {
        return errors.New("tanggal_selesai_kerja tidak boleh lebih awal dari tanggal_mulai_kerja")
    }
    
    return nil
}
