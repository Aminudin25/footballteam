package match_result

type CreateGoalInput struct {
	PlayerID int `json:"player_id" binding:"required"`
	TeamID   int `json:"team_id" binding:"required"`
	Minute   int `json:"minute" binding:"required"`
}

type CreateMatchResultInput struct {
	MatchID   int               `json:"match_id" binding:"required"`
	HomeScore int               `json:"home_score"`
	AwayScore int               `json:"away_score"`
	Status    string            `json:"status"`
	Goals     []CreateGoalInput `json:"goals"`
}
