// ...existing code...
package service

import (
    "errors"
    "time"

    "hello-fiber/app/model"
    "hello-fiber/app/repository"

    "github.com/gofiber/fiber/v2"
)

var alumniRepo = repository.NewAlumniRepositoryMongo()

// GetAllAlumniService mengambil semua alumni
func GetAllAlumniService(c *fiber.Ctx) error {
    alumni, err := alumniRepo.GetAllAlumni()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data alumni", "error": err.Error()})
    }
    return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data alumni", "data": alumni})
}

// GetAlumniByIDService mengambil alumni berdasarkan id (string)
func GetAlumniByIDService(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
    }

    item, err := alumniRepo.GetAlumniByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan", "error": err.Error()})
    }
    return c.JSON(fiber.Map{"success": true, "message": "Berhasil mengambil data alumni", "data": item})
}

// CreateAlumniService menambah alumni baru
func CreateAlumniService(c *fiber.Ctx) error {
    var req model.Alumni
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
    }

    // validasi minimal
    if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "NIM, nama, jurusan, dan email harus diisi"})
    }

    // set waktu (repository juga bisa set, tapi aman untuk set di sini)
    req.CreatedAt = time.Now()
    req.UpdatedAt = time.Now()

    id, err := alumniRepo.CreateAlumni(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menambahkan alumni", "error": err.Error()})
    }
    return c.Status(201).JSON(fiber.Map{"success": true, "message": "Alumni berhasil dibuat", "id": id})
}

// UpdateAlumniService mengupdate alumni berdasarkan id
func UpdateAlumniService(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
    }

    var req model.Alumni
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "Request body tidak valid", "error": err.Error()})
    }

    // validasi minimal (sesuaikan kebutuhan)
    if req.Nama == "" || req.Jurusan == "" || req.Email == "" {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "Nama, jurusan, dan email harus diisi"})
    }

    req.UpdatedAt = time.Now()
    if err := alumniRepo.UpdateAlumni(id, req); err != nil {
        return c.Status(404).JSON(fiber.Map{"success": false, "message": "Gagal mengupdate alumni", "error": err.Error()})
    }
    return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil diupdate"})
}

// DeleteAlumniService menghapus alumni berdasarkan id
func DeleteAlumniService(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"success": false, "message": "ID tidak boleh kosong"})
    }

    if err := alumniRepo.DeleteAlumni(id); err != nil {
        // jika repository mengembalikan error khusus, kembalikan 404
        if errors.Is(err, err) { // placeholder: repository mengembalikan error string biasa
            return c.Status(404).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan", "error": err.Error()})
        }
        return c.Status(500).JSON(fiber.Map{"success": false, "message": "Gagal menghapus alumni", "error": err.Error()})
    }
    return c.JSON(fiber.Map{"success": true, "message": "Alumni berhasil dihapus"})
}
// ...existing code...