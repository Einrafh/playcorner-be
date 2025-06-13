// Package routes
package routes

import (
	"playcorner-be/internal/handlers"
	"playcorner-be/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes menginisialisasi semua rute untuk aplikasi PlayCorner
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- Rute Publik ---
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh", handlers.RefreshToken)

	api.Get("/tvs", handlers.GetAllTVs)
	api.Get("/tvs/:tvId/reservations", handlers.GetTVReservations)

	// --- Rute Terproteksi ---
	// Rute di bawah ini memerlukan token JWT yang valid di header 'Authorization'
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.Get("/users/:userId", handlers.GetUser)
	protected.Get("/users/:userId/histories", handlers.GetUserHistories)
	protected.Post("/tvs/:tvId/reservations", handlers.CreateReservation)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "Welcome to PlayCorner API!"})
	})
}
