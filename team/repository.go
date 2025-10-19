package team

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Team, error)
	FindByID(id int) (Team, error)
	Create(team Team) (Team, error)
	Update(team Team) (Team, error)
	Delete(team Team) error
	FindByName(name string) (Team, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Team, error) {
	var teams []Team
	err := r.db.Find(&teams).Error
	return teams, err
}

func (r *repository) FindByID(id int) (Team, error) {
	var team Team
	err := r.db.First(&team, id).Error
	return team, err
}

func (r *repository) Create(team Team) (Team, error) {
	err := r.db.Create(&team).Error
	return team, err
}

func (r *repository) Update(team Team) (Team, error) {
	err := r.db.Save(&team).Error
	return team, err
}

func (r *repository) Delete(team Team) error {
	return r.db.Delete(&team).Error
}

func (r *repository) FindByName(name string) (Team, error) {
	var team Team
	err := r.db.Where("name = ?", name).First(&team).Error
	return team, err
}
