package player

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Height    float64   `gorm:"not null"`
	Weight    float64   `gorm:"not null"`
	Position  string    `gorm:"type:varchar(50);not null"` // Penyerang, Gelandang, Bertahan, Penjaga Gawang
	Number    int       `gorm:"not null"`
	TeamID    int       `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}