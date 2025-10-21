package route

import (
	"hello-fiber/app/service"
	"github.com/gofiber/fiber/v2"
	"hello-fiber/middleware"
)

func SetupRoutes(app *fiber.App, db interface{}) {
	api := app.Group("/api")

	// Route untuk registrasi dan login
	api.Post("/register", func(c *fiber.Ctx) error {
		return service.Register(c)
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c)
	})

	protected := api.Group("/", middleware.JWTMiddleware())

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
	pekerjaan.Delete("/trash/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.HardDeleteTrashedPekerjaanAlumniService(c)
	})
	pekerjaan.Put("/trash/:id/restore", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.RestoreTrashedPekerjaanAlumniService(c)
	})
	
	pekerjaan.Post("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c)
	})
	pekerjaan.Put("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c)
	})
	pekerjaan.Put("/admin/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniAdmin(c)
	})
	pekerjaan.Put("/users/:id", middleware.PekerjaanOwnerMiddlewareMongo(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniSementara(c)
	})
	pekerjaan.Delete("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c)
	})
}
