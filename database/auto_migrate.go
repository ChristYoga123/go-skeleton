package database

import (
	"golang-skeleton/app/entities/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		// ... tambahkan model lainnya di sini
	)
}
