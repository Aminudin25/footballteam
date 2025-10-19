package player

type PlayerFormatter struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Height   float64 `json:"height"`
	Weight   float64 `json:"weight"`
	Position string  `json:"position"`
	Number   int     `json:"number"`
	TeamID   int     `json:"team_id"`
}

func FormatPlayer(player Player) PlayerFormatter {
	return PlayerFormatter{
		ID:       player.ID,
		Name:     player.Name,
		Height:   player.Height,
		Weight:   player.Weight,
		Position: player.Position,
		Number:   player.Number,
		TeamID:   player.TeamID,
	}
}

func FormatPlayers(players []Player) []PlayerFormatter {
	formatted := []PlayerFormatter{}
	for _, p := range players {
		formatted = append(formatted, FormatPlayer(p))
	}
	return formatted
}
