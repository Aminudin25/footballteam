package match

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindBySchedule(date, time string, homeTeamID, awayTeamID int) (Match, error)
	FindAll() ([]Match, error)
	FindByID(id int) (Match, error)
	Create(match Match) (Match, error)
	Update(match Match) (Match, error)
	Delete(match Match) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindBySchedule(date, time string, homeTeamID, awayTeamID int) (Match, error) {
	var match Match
	err := r.db.
		Where("date = ? AND time = ? AND home_team_id = ? AND away_team_id = ?", date, time, homeTeamID, awayTeamID).
		First(&match).Error
	return match, err
}

func (r *repository) FindAll() ([]Match, error) {
	var matches []Match
	err := r.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Find(&matches).Error
	return matches, err
}

func (r *repository) FindByID(id int) (Match, error) {
	var match Match
	err := r.db.
		Preload("HomeTeam").
		Preload("AwayTeam").
		First(&match, id).Error
	return match, err
}

func (r *repository) Create(match Match) (Match, error) {
	err := r.db.Create(&match).Error
	return match, err
}

func (r *repository) Update(match Match) (Match, error) {
	err := r.db.Save(&match).Error
	return match, err
}

func (r *repository) Delete(match Match) error {
	// Soft delete otomatis karena pakai gorm.Model (DeletedAt)
	return r.db.Delete(&match).Error
}
