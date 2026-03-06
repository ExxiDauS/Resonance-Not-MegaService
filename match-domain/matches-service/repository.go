package matchesservice

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateMatch(match *Match) error {
	return r.db.Create(match).Error
}

func (r *Repository) GetMatchesByUserID(userID uuid.UUID) ([]Match, error) {
	var matches []Match
	err := r.db.Where("user_a_id = ? OR user_b_id = ?", userID, userID).Find(&matches).Error
	return matches, err
}

func (r *Repository) GetMatchByID(matchID uuid.UUID) (*Match, error) {
	var match Match
	err := r.db.Where("match_id = ?", matchID).First(&match).Error
	if err == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}
	return &match, err
}

func (r *Repository) GetAllMatches() ([]Match, error) {
	var matches []Match
	err := r.db.Find(&matches).Error
	return matches, err
}

func (r *Repository) DeleteMatch(matchID uuid.UUID) error {
	result := r.db.Where("match_id = ?", matchID).Delete(&Match{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
