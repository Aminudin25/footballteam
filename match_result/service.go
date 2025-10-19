package match_result

import (
	"fmt"
	"footballteam/match"
	"footballteam/player"
	"time"

	"gorm.io/gorm"
)

type Service interface {
	Create(input CreateMatchResultInput) (MatchResult, error)
	FindAll() ([]MatchResult, error)
	FindByID(id int) (MatchResult, error)
	GetMatchResultsReport() ([]MatchResultReportFormatter, error)
}

type service struct {
	repository     Repository
	playerService  player.Service // <-- tambahkan ini
	matchService   match.Service  // optional, untuk validasi match exist
}

func NewService(repo Repository, pService player.Service, mService match.Service) Service {
	return &service{
		repository:    repo,
		playerService: pService,
		matchService:  mService,
	}
}

func (s *service) Create(input CreateMatchResultInput) (MatchResult, error) {
    // Cek apakah match exist
    _, err := s.matchService.FindByID(input.MatchID)
    if err != nil {
        return MatchResult{}, fmt.Errorf("match with ID %d not found", input.MatchID)
    }

    // Cek apakah match result untuk match yang sama sudah ada
    existing, err := s.repository.FindByMatchID(input.MatchID)
    if err != nil && err != gorm.ErrRecordNotFound {
        return MatchResult{}, err
    }
    if err == nil {
        return existing, fmt.Errorf("match result for match ID %d already exists", input.MatchID)
    }

    // Buat MatchResult baru
    matchResult := MatchResult{
        MatchID:   input.MatchID,
        HomeScore: input.HomeScore,
        AwayScore: input.AwayScore,
        Status:    input.Status,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

   // Validasi goals
    for _, g := range input.Goals {
        player, err := s.playerService.GetPlayerByID(g.PlayerID)
        if err != nil {
            return MatchResult{}, fmt.Errorf("player with ID %d not found", g.PlayerID)
        }

        // Cek apakah player berada di tim yang sesuai
        if player.TeamID != g.TeamID {
            return MatchResult{}, fmt.Errorf("player %s (ID %d) does not belong to team %d", player.Name, player.ID, g.TeamID)
        }

        // Jika valid, masukkan goal
        matchResult.Goals = append(matchResult.Goals, Goal{
            PlayerID:  g.PlayerID,
            TeamID:    g.TeamID,
            Minute:    g.Minute,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        })
    }

    result, err := s.repository.Create(matchResult)
    if err != nil {
        return result, err
    }

    return result, nil
}


func (s *service) FindAll() ([]MatchResult, error) {
	return s.repository.FindAll()
}

func (s *service) FindByID(id int) (MatchResult, error) {
	return s.repository.FindByID(id)
}

func (s *service) GetMatchResultsReport() ([]MatchResultReportFormatter, error) {
	// Ambil semua match result beserta relasi Match, Teams, dan Goals
	results, err := s.repository.FindAllWithRelations()
	if err != nil {
		return nil, err
	}

	report := FormatMatchResultReport(results)
	return report, nil
}

