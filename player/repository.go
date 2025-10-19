package player

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Player, error)
	FindByID(id int) (Player, error)
	FindByTeamID(teamID int) ([]Player, error)
	Create(player Player) (Player, error)
	Update(player Player) (Player, error)
	Delete(player Player) error
	IsNumberExistInTeam(teamID, number int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Player, error) {
	var players []Player
	err := r.db.Find(&players).Error
	return players, err
}

func (r *repository) FindByID(id int) (Player, error) {
	var player Player
	err := r.db.First(&player, id).Error
	return player, err
}

func (r *repository) FindByTeamID(teamID int) ([]Player, error) {
	var players []Player
	err := r.db.Where("team_id = ?", teamID).Find(&players).Error
	return players, err
}

func (r *repository) Create(player Player) (Player, error) {
	err := r.db.Create(&player).Error
	return player, err
}

func (r *repository) Update(player Player) (Player, error) {
	err := r.db.Save(&player).Error
	return player, err
}

func (r *repository) Delete(player Player) error {
	return r.db.Delete(&player).Error
}

// Cek apakah nomor punggung sudah digunakan dalam satu tim
func (r *repository) IsNumberExistInTeam(teamID, number int) (bool, error) {
	var count int64
	err := r.db.Model(&Player{}).
		Where("team_id = ? AND number = ?", teamID, number).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
