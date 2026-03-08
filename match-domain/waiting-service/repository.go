package waitingservice

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

func (r *Repository) FindByTrackID(trackID string) (*Waiting, error) {
	var waiting Waiting
	err := r.db.Where("track_id = ?", trackID).First(&waiting).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &waiting, err
}

func (r *Repository) AddToWaiting(waiting *Waiting) error {
	return r.db.Create(waiting).Error
}

func (r *Repository) DeleteByID(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&Waiting{}).Error
}
