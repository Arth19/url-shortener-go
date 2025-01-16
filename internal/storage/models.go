package storage

import (
	"time"
)

type URL struct {
	ID         uint   `gorm:"primaryKey"`
	ShortCode  string `gorm:"uniqueIndex"`
	Original   string
	ClickCount uint `gorm:"default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
