// main.go
package main

import (
	"log"
	"os"
	"playcorner-be/internal/database"
	"playcorner-be/internal/models"
	"playcorner-be/internal/routes"
	"playcorner-be/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	database.ConnectDB()

	// Menambahkan data awal ke database jika belum ada
	seedDatabase()

	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}

// seedDatabase membuat data awal untuk testing
func seedDatabase() {
	// Cek apakah data sudah ada untuk menghindari duplikasi
	var tvCount int64
	database.DB.Model(&models.TVInfo{}).Count(&tvCount)

	if tvCount == 0 {
		log.Println("Seeding database with fixed relationships...")

		// Buat User
		hashedPassword, _ := utils.HashPassword("password123")
		user := models.User{
			ID:             "235150207111062",
			Name:           "Muhammad Rafly Ash Shiddiqi",
			Faculty:        "FILKOM",
			Major:          "Teknik Informatika",
			CreditScore:    100,
			ProfilePictURL: "https://i.pravatar.cc/150?u=235150207111062",
			PasswordHash:   hashedPassword,
		}
		database.DB.Create(&user)

		// Definisikan semua game sebagai pointer
		gameFC24 := &models.Game{Title: "EA Sports FC 24", CoverPictURL: "https://placehold.co/200x300/3498DB/FFFFFF?text=FC+24"}
		gameTekken8 := &models.Game{Title: "Tekken 8", CoverPictURL: "https://placehold.co/200x300/E74C3C/FFFFFF?text=Tekken+8"}
		gameItTakesTwo := &models.Game{Title: "It Takes Two", CoverPictURL: "https://placehold.co/200x300/F1C40F/FFFFFF?text=It+Takes+Two"}
		gameOvercooked := &models.Game{Title: "Overcooked! 2", CoverPictURL: "https://placehold.co/200x300/9B59B6/FFFFFF?text=Overcooked"}

		// Buat semua game dari slice of pointers. GORM akan mengisi ID kembali ke variabel asli.
		allGames := []*models.Game{gameFC24, gameTekken8, gameItTakesTwo, gameOvercooked}
		if err := database.DB.Create(allGames).Error; err != nil {
			log.Fatalf("Failed to seed games: %v", err)
		}

		// 3. Definisikan dan buat TV, sekarang mereferensikan game yang sudah memiliki ID
		tvsToCreate := []models.TVInfo{
			{
				ID:          1,
				ConsoleType: "PlayStation 5",
				Games:       []*models.Game{gameFC24, gameTekken8}, // Sekarang gameFC24.ID bukan lagi 0
			},
			{
				ID:          2,
				ConsoleType: "PlayStation 5",
				Games:       []*models.Game{gameFC24, gameItTakesTwo},
			},
			{
				ID:          3,
				ConsoleType: "Xbox Series X",
				Games:       []*models.Game{gameFC24, gameOvercooked},
			},
			{
				ID:          4,
				ConsoleType: "Xbox Series X",
				Games:       []*models.Game{gameFC24},
			},
			{
				ID:          5,
				ConsoleType: "PC",
				Games:       []*models.Game{gameFC24, gameItTakesTwo, gameTekken8},
			},
		}

		if err := database.DB.Create(&tvsToCreate).Error; err != nil {
			log.Fatalf("Failed to seed tvs with associations: %v", err)
		}

		log.Println("Seeding complete.")
	}
}
