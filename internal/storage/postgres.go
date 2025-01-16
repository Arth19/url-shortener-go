package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=url_shortener port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("erro ao conectar ao Postgres: %w", err)
	}

	err = db.AutoMigrate(&URL{})
	if err != nil {
		return fmt.Errorf("erro no auto-migrate: %w", err)
	}

	DB = db
	log.Println("Conexão com Postgres iniciada com sucesso!")
	return nil
}

func SaveURL(shortCode, original string) (*URL, error) {
	url := &URL{
		ShortCode: shortCode,
		Original:  original,
	}
	if err := DB.Create(url).Error; err != nil {
		return nil, err
	}
	return url, nil
}

func GetURL(shortCode string) (*URL, error) {
	var url URL
	if err := DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return nil, err
	}
	return &url, nil
}

func IncrementClickCount(shortCode string) error {
	return DB.Model(&URL{}).
		Where("short_code = ?", shortCode).
		Update("click_count", gorm.Expr("click_count + ?", 1)).
		Error
}
