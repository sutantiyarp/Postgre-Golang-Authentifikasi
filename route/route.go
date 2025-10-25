package route

import (
	"hello-fiber/app/repository"
	"hello-fiber/app/service"
	"github.com/gofiber/fiber/v2"
	"hello-fiber/middleware"
)

func SetupRoutes(app *fiber.App, db interface{}) {
	api := app.Group("/api")

	api.Post("/register", func(c *fiber.Ctx) error {
		return service.Register(c)
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c)
	})

	protected := api.Group("/", middleware.JWTMiddleware())

	users := protected.Group("/users")
	users.Get("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.GetAllUsersService(c)
	})
	users.Post("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.CreateUserAdmin(c)
	})
	users.Put("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdateUserService(c)
	})
	users.Delete("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.DeleteUserService(c)
	})

	alumni := protected.Group("/alumni")
	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c)
	})
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c)
	})
	alumni.Post("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c)
	})
	alumni.Put("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c)
	})
	alumni.Delete("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c)
	})

	pekerjaan := protected.Group("/pekerjaan")
	
	// Get routes
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c)
	})
	pekerjaan.Get("/trash", func(c *fiber.Ctx) error {
		return service.GetTrashedPekerjaanAlumniService(c)
	})
	pekerjaan.Get("/alumni/:alumni_id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByAlumniIDService(c)
	})
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByIDService(c)
	})
	
	// Create route
	pekerjaan.Post("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c)
	})
	
	// Update routes
	pekerjaan.Put("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c)
	})
	
	// Delete routes (soft delete)
	pekerjaan.Delete("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c)
	})
	
	// Trash management routes
	pekerjaan.Delete("/trash/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.HardDeleteTrashedPekerjaanAlumniService(c)
	})
	pekerjaan.Put("/trash/:id/restore", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.RestoreTrashedPekerjaanAlumniService(c)
	})

	fileUploadRepo := repository.NewFileUploadRepository()
	fileUploadService := service.NewFileUploadService(fileUploadRepo, "./uploads")

	files := protected.Group("/files")
	files.Get("/", func(c *fiber.Ctx) error {
		return fileUploadService.GetAllFiles(c)
	})
	files.Get("/:id", func(c *fiber.Ctx) error {
		return fileUploadService.GetFileByID(c)
	})
	files.Post("/foto", func(c *fiber.Ctx) error {
		return fileUploadService.UploadFoto(c)
	})
	files.Post("/sertifikat", func(c *fiber.Ctx) error {
		return fileUploadService.UploadSertifikat(c)
	})
	files.Delete("/:id", func(c *fiber.Ctx) error {
		return fileUploadService.DeleteFile(c)
	})
}