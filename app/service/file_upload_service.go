package service

import (
	"fmt"
	"os"
	"path/filepath"

	"hello-fiber/app/model"
	"hello-fiber/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type FileUploadService interface {
	UploadFoto(c *fiber.Ctx) error
	UploadSertifikat(c *fiber.Ctx) error
	GetAllFiles(c *fiber.Ctx) error
	GetFileByID(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type fileUploadService struct {
	repo       repository.FileUploadRepository
	uploadPath string
}

func NewFileUploadService(repo repository.FileUploadRepository, uploadPath string) FileUploadService {
	return &fileUploadService{
		repo:       repo,
		uploadPath: uploadPath,
	}
}

func (s *fileUploadService) UploadFoto(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	// Validasi ukuran file (max 1MB)
	maxSize := int64(1 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File size exceeds 1MB limit",
		})
	}

	// Validasi tipe file
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File type not allowed. Only JPEG, JPG, and PNG are allowed",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := "foto_" + uuid.New().String() + ext
	uploadDir := filepath.Join(s.uploadPath, "foto")
	filePath := filepath.Join(uploadDir, newFileName)

	// Create directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	// Save file
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}

	fileModel := &model.FileUpload{
		FileType:     "foto",
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file metadata",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Foto uploaded successfully",
		"data":    s.toFileResponse(fileModel),
	})
}

func (s *fileUploadService) UploadSertifikat(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "No file uploaded",
			"error":   err.Error(),
		})
	}

	// Validasi ukuran file (max 2MB)
	maxSize := int64(2 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File size exceeds 2MB limit",
		})
	}

	// Validasi tipe file
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File type not allowed. Only PDF is allowed",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := "sertifikat_" + uuid.New().String() + ext
	uploadDir := filepath.Join(s.uploadPath, "sertifikat")
	filePath := filepath.Join(uploadDir, newFileName)

	// Create directory if not exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
	}

	// Save file
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
	}

	fileModel := &model.FileUpload{
		FileType:     "sertifikat",
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
	}

	if err := s.repo.Create(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to save file metadata",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Sertifikat uploaded successfully",
		"data":    s.toFileResponse(fileModel),
	})
}

func (s *fileUploadService) GetAllFiles(c *fiber.Ctx) error {
	files, err := s.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get files",
			"error":   err.Error(),
		})
	}

	var responses []model.FileUploadResponse
	for _, file := range files {
		responses = append(responses, *s.toFileResponse(&file))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Files retrieved successfully",
		"data":    responses,
	})
}

func (s *fileUploadService) GetFileByID(c *fiber.Ctx) error {
	fileIDStr := c.Params("id")
	if fileIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File ID is required",
		})
	}

	fileID, err := bson.ObjectIDFromHex(fileIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid file ID format",
		})
	}

	file, err := s.repo.FindByID(fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File retrieved successfully",
		"data":    s.toFileResponse(file),
	})
}

func (s *fileUploadService) DeleteFile(c *fiber.Ctx) error {
	fileIDStr := c.Params("id")
	if fileIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "File ID is required",
		})
	}

	fileID, err := bson.ObjectIDFromHex(fileIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid file ID format",
		})
	}

	file, err := s.repo.FindByID(fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File not found",
		})
	}

	// Delete file from storage
	if err := os.Remove(file.FilePath); err != nil {
		fmt.Println("Warning: Failed to delete file from storage:", err)
	}

	// Delete from database
	if err := s.repo.Delete(fileID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete file",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "File deleted successfully",
	})
}

func (s *fileUploadService) toFileResponse(file *model.FileUpload) *model.FileUploadResponse {
	return &model.FileUploadResponse{
		ID:           file.ID,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		UploadedAt:   file.UploadedAt,
	}
}
