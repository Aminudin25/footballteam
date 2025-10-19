package match_result

import "gorm.io/gorm"

type Repository interface {
	Create(result MatchResult) (MatchResult, error)
	FindByID(id int) (MatchResult, error)
	FindAll() ([]MatchResult, error)
	FindByMatchID(matchID int) (MatchResult, error)
	FindAllWithRelations() ([]MatchResult, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByMatchID(matchID int) (MatchResult, error) {
	var matchResult MatchResult
	err := r.db.Preload("Goals").Where("match_id = ?", matchID).First(&matchResult).Error
	return matchResult, err
}

func (r *repository) Create(matchResult MatchResult) (MatchResult, error) {
	err := r.db.Create(&matchResult).Error
	return matchResult, err
}


func (r *repository) FindByID(id int) (MatchResult, error) {
	var result MatchResult
	err := r.db.Preload("Goals.Player").First(&result, id).Error
	return result, err
}

func (r *repository) FindAll() ([]MatchResult, error) {
	var results []MatchResult
	err := r.db.Preload("Goals.Player").Find(&results).Error
	return results, err
}

func (r *repository) FindAllWithRelations() ([]MatchResult, error) {
    var results []MatchResult
    err := r.db.Preload("Goals.Player").
               Preload("Match.HomeTeam").
               Preload("Match.AwayTeam").
               Find(&results).Error
    return results, err
}

