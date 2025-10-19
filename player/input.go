package player

type CreatePlayerInput struct {
	Name     string  `json:"name" binding:"required"`
	Height   float64 `json:"height" binding:"required"`
	Weight   float64 `json:"weight" binding:"required"`
	Position string  `json:"position" binding:"required"`
	Number   int     `json:"number" binding:"required"`
	TeamID   int     `json:"team_id" binding:"required"`
}

type UpdatePlayerInput struct {
	Name     string  `json:"name"`
	Height   float64 `json:"height"`
	Weight   float64 `json:"weight"`
	Position string  `json:"position"`
	Number   int     `json:"number"`
	TeamID   int     `json:"team_id"`
}
