package match

type TeamFormatter struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MatchFormatter struct {
	ID       int           `json:"id"`
	Date     string        `json:"date"`
	Time     string        `json:"time"`
	HomeTeam TeamFormatter `json:"home_team"`
	AwayTeam TeamFormatter `json:"away_team"`
}

func FormatMatch(m Match) MatchFormatter {
	return MatchFormatter{
		ID:   m.ID,
		Date: m.Date,
		Time: m.Time,
		HomeTeam: TeamFormatter{
			ID:   m.HomeTeam.ID,
			Name: m.HomeTeam.Name,
		},
		AwayTeam: TeamFormatter{
			ID:   m.AwayTeam.ID,
			Name: m.AwayTeam.Name,
		},
	}
}

// 🔥 Tambahan untuk list
func FormatMatches(matches []Match) []MatchFormatter {
	formatted := []MatchFormatter{}
	for _, m := range matches {
		formatted = append(formatted, FormatMatch(m))
	}
	return formatted
}
