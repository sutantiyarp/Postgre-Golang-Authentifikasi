package service

import (
	"errors"
	"time"

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

	// validasi minimal (sesuaikan dengan kebutuhan)
	if err := validatePekerjaanInput(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	// Convert AlumniID string to bson.ObjectID
	alumniID, err := bson.ObjectIDFromHex(req.AlumniID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "alumni_id format tidak valid"})
	}

	pekerjaan := model.PekerjaanAlumni{
		AlumniID:            alumniID,
		NamaPerusahaan:      req.NamaPerusahaan,
		PosisiJabatan:       req.PosisiJabatan,
		BidangIndustri:      req.BidangIndustri,
		LokasiKerja:         req.LokasiKerja,
		GajiRange:           req.GajiRange,
		TanggalMulaiKerja:   req.TanggalMulaiKerja,
		TanggalSelesaiKerja: req.TanggalSelesaiKerja,
		StatusPekerjaan:     req.StatusPekerjaan,
		DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
		IsDelete:            "tidak",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if pekerjaan.StatusPekerjaan == "" {
		pekerjaan.StatusPekerjaan = "aktif"
	}

	id, err := pekerjaanRepo.CreatePekerjaanAlumni(pekerjaan)
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

	var req model.PekerjaanAlumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	req.UpdatedAt = time.Now()
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

	// body optional: {"is_delete":"hapus"} atau {"is_delete":"tidak"}
	var body struct {
		IsDelete string `json:"is_delete"`
	}
	if err := c.BodyParser(&body); err != nil {
		body.IsDelete = "hapus"
	}
	if body.IsDelete == "" {
		body.IsDelete = "hapus"
	}

	if err := pekerjaanRepo.SoftDeletePekerjaanAlumni(id, body.IsDelete); err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengubah status is_delete", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Status is_delete berhasil diupdate"})
}

// GetTrashedPekerjaanAlumniService mengambil semua yang is_delete == "hapus"
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

// UpdatePekerjaanAlumniAdmin
func UpdatePekerjaanAlumniAdmin(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	var req model.PekerjaanAlumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	req.UpdatedAt = time.Now()
	err = pekerjaanRepo.UpdatePekerjaanAlumni(id, req)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengupdate pekerjaan alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan alumni berhasil diupdate"})
}

// UpdatePekerjaanAlumniSementara
func UpdatePekerjaanAlumniSementara(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	var req model.PekerjaanAlumni
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
	}

	req.UpdatedAt = time.Now()
	err = pekerjaanRepo.UpdatePekerjaanAlumni(id, req)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengupdate pekerjaan alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan alumni berhasil diupdate"})
}

// DeletePekerjaanAlumniService
func DeletePekerjaanAlumniService(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
	}

	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID format tidak valid"})
	}

	if err := pekerjaanRepo.SoftDeletePekerjaanAlumni(id, "hapus"); err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal menghapus pekerjaan alumni", "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan alumni berhasil dihapus"})
}

// helper untuk validasi minimal (opsional)
func validatePekerjaanInput(p model.CreatePekerjaanAlumniRequest) error {
	if p.AlumniID == "" {
		return errors.New("alumni_id wajib diisi")
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
	// TanggalMulaiKerja bertipe time.Time
	if p.TanggalMulaiKerja.IsZero() {
		return errors.New("tanggal_mulai_kerja wajib diisi")
	}
	return nil
}
