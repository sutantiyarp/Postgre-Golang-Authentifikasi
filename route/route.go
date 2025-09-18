package route

import (
	"hello-fiber/app/service"
	"database/sql"
	"github.com/gofiber/fiber/v2"
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

	alumni := api.Group("/alumni")
	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})
	alumni.Post("/", func(c *fiber.Ctx) error {
		return service.CreateAlumniService(c, db)
	})
	alumni.Put("/:id", func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})
	alumni.Delete("/:id", func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	pekerjaan := api.Group("/pekerjaan")
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c, db)
	})
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByIDService(c, db)
	})
	pekerjaan.Get("/alumni/:alumni_id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanAlumniByAlumniIDService(c, db)
	})
	pekerjaan.Post("/", func(c *fiber.Ctx) error {
		return service.CreatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Put("/:id", func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c, db)
	})
	pekerjaan.Delete("/:id", func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c, db)
	})
}