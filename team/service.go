package team

import "fmt"

type Service interface {
	GetAllTeams() ([]Team, error)
	GetTeamByID(id int) (Team, error)
	CreateTeam(input CreateTeamInput) (Team, error)
	UpdateTeam(id int, input UpdateTeamInput) (Team, error)
	DeleteTeam(id int) error
	SaveLogo(id int, fileLocation string) (Team, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetAllTeams() ([]Team, error) {
	return s.repository.FindAll()
}

func (s *service) GetTeamByID(id int) (Team, error) {
	return s.repository.FindByID(id)
}

func (s *service) CreateTeam(input CreateTeamInput) (Team, error) {
	team := Team{
		Name:        input.Name,
		YearFounded: input.YearFounded,
		Address:     input.Address,
		City:        input.City,
	}
	return s.repository.Create(team)
}

func (s *service) UpdateTeam(id int, input UpdateTeamInput) (Team, error) {
	team, err := s.repository.FindByID(id)
	if err != nil {
		return team, err
	}

	if input.Name != "" {
		team.Name = input.Name
	}
	if input.YearFounded != 0 {
		team.YearFounded = input.YearFounded
	}
	if input.Address != "" {
		team.Address = input.Address
	}
	if input.City != "" {
		team.City = input.City
	}

	return s.repository.Update(team)
}

func (s *service) DeleteTeam(id int) error {
	team, err := s.repository.FindByID(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(team)
}

func (s *service) SaveLogo(id int, fileLocation string) (Team, error) {
	fmt.Println("🔍 SaveLogo() called with ID:", id)
	team, err := s.repository.FindByID(id)
	fmt.Println("🔍 Team result:", team, "Error:", err)
	if err != nil {
		return team, err
	}

	team.Logo = fileLocation

	updateTeam, err := s.repository.Update(team)

	if err != nil {
		return updateTeam, err
	}

	return updateTeam, nil
}
