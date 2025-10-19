package match_result

import (
	"sort"
)

// Formatter untuk Goal per MatchResult
type GoalFormatter struct {
	PlayerID   int    `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamID     int    `json:"team_id"`
	Minute     int    `json:"minute"`
}

// Formatter untuk response MatchResult biasa
type MatchResultFormatter struct {
	ID        int             `json:"id"`
	MatchID   int             `json:"match_id"`
	HomeScore int             `json:"home_score"`
	AwayScore int             `json:"away_score"`
	Status    string          `json:"status"`
	Goals     []GoalFormatter `json:"goals"`
}

// Formatter untuk response report
type MatchResultReportFormatter struct {
	MatchID       int      `json:"match_id"`
	Date          string   `json:"date"`
	Time          string   `json:"time"`
	HomeTeam      string   `json:"home_team"`
	AwayTeam      string   `json:"away_team"`
	HomeScore     int      `json:"home_score"`
	AwayScore     int      `json:"away_score"`
	Status        string   `json:"status"`          // Home Menang / Away Menang / Draw
	TopScorer     []string `json:"top_scorers"`     // nama pemain dengan gol terbanyak
	HomeTotalWins int      `json:"home_total_wins"` // akumulasi kemenangan Home Team hingga match ini
	AwayTotalWins int      `json:"away_total_wins"` // akumulasi kemenangan Away Team hingga match ini
}

// FormatMatchResult untuk response biasa
func FormatMatchResult(m MatchResult, playerMap map[int]string) MatchResultFormatter {
	var goals []GoalFormatter
	for _, g := range m.Goals {
		goals = append(goals, GoalFormatter{
			PlayerID:   g.PlayerID,
			PlayerName: playerMap[g.PlayerID],
			TeamID:     g.TeamID,
			Minute:     g.Minute,
		})
	}

	return MatchResultFormatter{
		ID:        m.ID,
		MatchID:   m.MatchID,
		HomeScore: m.HomeScore,
		AwayScore: m.AwayScore,
		Status:    m.Status,
		Goals:     goals,
	}
}

// FormatMatchResultReport untuk report lengkap
func FormatMatchResultReport(results []MatchResult) []MatchResultReportFormatter {
	report := []MatchResultReportFormatter{}

	// total kemenangan per tim
	homeWinsMap := make(map[int]int)
	awayWinsMap := make(map[int]int)

	for _, r := range results {
		topScorerMap := make(map[string]int)

		// hitung gol per pemain
		for _, g := range r.Goals {
			playerName := g.Player.Name
			topScorerMap[playerName]++
		}

		// cari pemain dengan gol terbanyak
		topScorers := []string{}
		maxGoals := 0
		for name, count := range topScorerMap {
			if count > maxGoals {
				maxGoals = count
				topScorers = []string{name}
			} else if count == maxGoals {
				topScorers = append(topScorers, name)
			}
		}

		// akumulasi kemenangan
		switch r.Status {
		case "Home Menang":
			homeWinsMap[r.MatchID]++
		case "Away Menang":
			awayWinsMap[r.MatchID]++
		case "Draw":
			// tidak menambah kemenangan
		default:
			// bisa log atau ignore
		}

		m := r.Match // pastikan MatchResult memiliki relasi Match

		report = append(report, MatchResultReportFormatter{
			MatchID:       r.MatchID,
			Date:          m.Date,
			Time:          m.Time,
			HomeTeam:      m.HomeTeam.Name,
			AwayTeam:      m.AwayTeam.Name,
			HomeScore:     r.HomeScore,
			AwayScore:     r.AwayScore,
			Status:        r.Status,
			TopScorer:     topScorers,
			HomeTotalWins: homeWinsMap[r.MatchID],
			AwayTotalWins: awayWinsMap[r.MatchID],
		})
	}

	// urutkan berdasarkan tanggal & waktu
	sort.Slice(report, func(i, j int) bool {
		if report[i].Date == report[j].Date {
			return report[i].Time < report[j].Time
		}
		return report[i].Date < report[j].Date
	})

	return report
}
