package team

type CreateTeamInput struct {
	Name        string `json:"name" binding:"required"`
	YearFounded int    `json:"year_founded" binding:"required"`
	Address     string `json:"address"`
	City        string `json:"city"`
}

type UpdateTeamInput struct {
	Name        string `json:"name"`
	YearFounded int    `json:"year_founded"`
	Address     string `json:"address"`
	City        string `json:"city"`
}
