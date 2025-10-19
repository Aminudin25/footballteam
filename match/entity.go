package match

import (
	"footballteam/team"
	"time"
)

type Match struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Date        string     `json:"date"`
	Time        string     `json:"time"`
	HomeTeamID  int        `json:"home_team_id"`
	AwayTeamID  int        `json:"away_team_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	HomeTeam team.Team `gorm:"foreignKey:HomeTeamID"`
	AwayTeam team.Team `gorm:"foreignKey:AwayTeamID"`
}
