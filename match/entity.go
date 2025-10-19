package match

import (
	"footballteam/team"
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID          int            `gorm:"primaryKey;autoIncrement"`
	Date        string         `gorm:"type:varchar(20);not null"`
	Time        string         `gorm:"type:varchar(10);not null"`
	HomeTeamID  int            `gorm:"not null"`
	AwayTeamID  int            `gorm:"not null"`
	HomeTeam   team.Team      `gorm:"foreignKey:HomeTeamID"`
	AwayTeam   team.Team      `gorm:"foreignKey:AwayTeamID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
