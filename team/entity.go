package team

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          	int            `gorm:"primaryKey;autoIncrement"`
	Name        	string         `gorm:"size:100;not null"`
	Logo        	string         `gorm:"size:255"`
	YearFounded 	int            `gorm:"not null"`
	Address     	string         `gorm:"size:255"`
	City        	string         `gorm:"size:100"`
	CreatedAt   	time.Time      `json:"created_at"`
	UpdatedAt   	time.Time      `json:"updated_at"`
	DeletedAt   	gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
