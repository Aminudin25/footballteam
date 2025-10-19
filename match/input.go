package match

type CreateMatchInput struct {
	Date       string `json:"date" binding:"required"`
	Time       string `json:"time" binding:"required"`
	HomeTeamID int    `json:"home_team_id" binding:"required"`
	AwayTeamID int    `json:"away_team_id" binding:"required"`
}

type UpdateMatchInput struct {
	Date       string `json:"date" binding:"required"`
	Time       string `json:"time" binding:"required"`
	HomeTeamID int    `json:"home_team_id" binding:"required"`
	AwayTeamID int    `json:"away_team_id" binding:"required"`
}
