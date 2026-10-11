package config

import (
	"log"
	"os"

	"even.com/even-app/models" // Sesuaikan dengan nama module/import project Anda

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	// Perbaiki dari model.Event menjadi models.Event
	err = DB.AutoMigrate(&models.Event{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}
	log.Println("Berhasil terhubung ke Neon DB!")
}
