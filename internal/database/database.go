// Package database
package database

import (
	"fmt"
	"log"
	"os"
	"playcorner-be/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah variabel global yang akan menampung instance koneksi GORM
// dan akan diakses oleh package lain (seperti handlers).
var DB *gorm.DB

// ConnectDB adalah satu-satunya fungsi yang perlu dipanggil dari main.go
// untuk menginisialisasi koneksi dan skema database.
func ConnectDB() {
	// Memuat variabel environment dari file .env.
	// Ini adalah praktik yang baik untuk menjaga kerahasiaan kredensial.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables from system")
	}

	// Membangun DSN (Data Source Name) dari environment variables.
	// Pastikan file .env Anda menggunakan nama variabel ini.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Membuka koneksi ke database menggunakan GORM.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Menampilkan query SQL di log untuk debugging.
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}
	log.Println("Database connection established")

	// Menjalankan AutoMigrate untuk membuat/memperbarui tabel database
	// secara otomatis sesuai dengan struct yang didefinisikan di package models.
	log.Println("Running Migrations")
	err = db.AutoMigrate(&models.User{}, &models.TVInfo{}, &models.Game{}, &models.Reservation{})
	if err != nil {
		log.Fatal("Migration failed. \n", err)
	}
	log.Println("Migrations completed")

	// Menetapkan instance database yang berhasil terhubung ke variabel global DB.
	DB = db
}
