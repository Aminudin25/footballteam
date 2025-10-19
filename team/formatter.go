package team

type TeamFormatter struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	YearFounded int    `json:"year_founded"`
	Address     string `json:"address"`
	City        string `json:"city"`
}

func FormatTeam(team Team) TeamFormatter {
	return TeamFormatter{
		ID:          team.ID,
		Name:        team.Name,
		Logo:        team.Logo,
		YearFounded: team.YearFounded,
		Address:     team.Address,
		City:        team.City,
	}
}
