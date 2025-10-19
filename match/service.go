package match

import "errors"

type Service interface {
	FindAll() ([]Match, error)
	FindByID(id int) (Match, error)
	CreateMatch(input CreateMatchInput) (Match, error)
	UpdateMatch(id int, input UpdateMatchInput) (Match, error)
	DeleteMatch(id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() ([]Match, error) {
	return s.repository.FindAll()
}

func (s *service) FindByID(id int) (Match, error) {
	return s.repository.FindByID(id)
}

func (s *service) CreateMatch(input CreateMatchInput) (Match, error) {
	// Cek apakah match sudah ada dengan jadwal yang sama
	existing, err := s.repository.FindBySchedule(input.Date, input.Time, input.HomeTeamID, input.AwayTeamID)
	if err == nil && existing.ID != 0 {
		return Match{}, errors.New("pertandingan dengan jadwal, jam, dan tim yang sama sudah terdaftar")
	}

	match := Match{
		Date:       input.Date,
		Time:       input.Time,
		HomeTeamID: input.HomeTeamID,
		AwayTeamID: input.AwayTeamID,
	}

	newMatch, err := s.repository.Create(match)
	if err != nil {
		return newMatch, err
	}

	return newMatch, nil
}

func (s *service) UpdateMatch(id int, input UpdateMatchInput) (Match, error) {
	match, err := s.repository.FindByID(id)
	if err != nil {
		return match, errors.New("match not found")
	}

	// Cek duplikasi (kecuali dirinya sendiri)
	existing, err := s.repository.FindBySchedule(input.Date, input.Time, input.HomeTeamID, input.AwayTeamID)
	if err == nil && existing.ID != 0 && existing.ID != id {
		return Match{}, errors.New("pertandingan dengan jadwal, jam, dan tim yang sama sudah terdaftar")
	}

	match.Date = input.Date
	match.Time = input.Time
	match.HomeTeamID = input.HomeTeamID
	match.AwayTeamID = input.AwayTeamID

	updated, err := s.repository.Update(match)
	if err != nil {
		return updated, err
	}

	return updated, nil
}

func (s *service) DeleteMatch(id int) error {
	match, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(match)
}
