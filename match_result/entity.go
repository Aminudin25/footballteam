package match_result

import (
	"time"

	"footballteam/match"
	"footballteam/player"
)

type MatchResult struct {
	ID        int        `gorm:"primaryKey" json:"id"`
	MatchID   int        `json:"match_id"`
	HomeScore int        `json:"home_score"`
	AwayScore int        `json:"away_score"`
	Status    string     `json:"status"` // Home Menang, Away Menang, Draw
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Match  match.Match `gorm:"foreignKey:MatchID"`  // untuk akses HomeTeam & AwayTeam
	Goals  []Goal      `gorm:"foreignKey:MatchResultID"`
}

type Goal struct {
	ID            int        `gorm:"primaryKey" json:"id"`
	MatchResultID int        `json:"match_result_id"`
	PlayerID      int        `json:"player_id"`
	TeamID        int        `json:"team_id"`
	Minute        int        `json:"minute"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Player player.Player `gorm:"foreignKey:PlayerID"`
}
