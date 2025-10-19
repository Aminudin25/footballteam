package player

import (
	"errors"
	"time"
)

type Service interface {
	GetAllPlayers() ([]Player, error)
	GetPlayerByID(id int) (Player, error)
	GetPlayersByTeam(teamID int) ([]Player, error)
	CreatePlayer(input CreatePlayerInput) (Player, error)
	UpdatePlayer(id int, input UpdatePlayerInput) (Player, error)
	DeletePlayer(id int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetAllPlayers() ([]Player, error) {
	return s.repository.FindAll()
}

func (s *service) GetPlayerByID(id int) (Player, error) {
	return s.repository.FindByID(id)
}

func (s *service) GetPlayersByTeam(teamID int) ([]Player, error) {
	return s.repository.FindByTeamID(teamID)
}

func (s *service) CreatePlayer(input CreatePlayerInput) (Player, error) {
	exist, err := s.repository.IsNumberExistInTeam(input.TeamID, input.Number)
	if err != nil {
		return Player{}, err
	}
	if exist {
		return Player{}, errors.New("nomor punggung sudah digunakan oleh pemain lain di tim ini")
	}

	player := Player{
		Name:      input.Name,
		Height:    input.Height,
		Weight:    input.Weight,
		Position:  input.Position,
		Number:    input.Number,
		TeamID:    input.TeamID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repository.Create(player)
}

func (s *service) UpdatePlayer(id int, input UpdatePlayerInput) (Player, error) {
	player, err := s.repository.FindByID(id)
	if err != nil {
		return player, err
	}

	// Jika mengganti nomor, cek duplikasi
	if input.Number != 0 && input.Number != player.Number {
		exist, err := s.repository.IsNumberExistInTeam(player.TeamID, input.Number)
		if err != nil {
			return player, err
		}
		if exist {
			return player, errors.New("nomor punggung sudah digunakan oleh pemain lain di tim ini")
		}
	}

	if input.Name != "" {
		player.Name = input.Name
	}
	if input.Height != 0 {
		player.Height = input.Height
	}
	if input.Weight != 0 {
		player.Weight = input.Weight
	}
	if input.Position != "" {
		player.Position = input.Position
	}
	if input.Number != 0 {
		player.Number = input.Number
	}
	if input.TeamID != 0 {
		player.TeamID = input.TeamID
	}
	player.UpdatedAt = time.Now()

	return s.repository.Update(player)
}

func (s *service) DeletePlayer(id int) error {
	player, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(player)
}
