package route

import (
	"hello-fiber/app/service"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"hello-fiber/middleware"
	_ "github.com/lib/pq" // <— WAJIB: daftarkan driver postgres
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	// Route untuk registrasi dan login
	api.Post("/register", func(c *fiber.Ctx) error {
		return service.RegisterService(c, db)
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

protected := api.Group("/", middleware.JWTMiddleware())

	alumni := protected.Group("/alumni")
	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})
	
	alumni.Post("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})
	alumni.Put("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})
	alumni.Delete("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	pekerjaan := protected.Group("/pekerjaan")
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c, db)
	})
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByIDService(c, db)
	})
	pekerjaan.Get("/alumni/:alumni_id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByAlumniIDService(c, db)
	})
	
	pekerjaan.Post("/", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("users/:id", middleware.PekerjaanOwnerMiddleware(db), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniSementara(c, db)
	})
	pekerjaan.Delete("/:id", middleware.AdminOnlyMiddleware(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c, db)
	})

	trash := pekerjaan.Group("/trash", middleware.AdminOnlyMiddleware())
	trash.Get("/", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniTrashService(c, db)
	})
	trash.Post("/:id/restore", func(c *fiber.Ctx) error {
		return service.RestorePekerjaanAlumniTrashService(c, db)
	})
	trash.Delete("/:id", func(c *fiber.Ctx) error {
		return service.HardDeletePekerjaanAlumniTrashService(c, db)
	})
}