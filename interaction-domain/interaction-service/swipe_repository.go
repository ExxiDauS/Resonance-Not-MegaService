package interactionservice

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type SwipeRepository interface {
	Save(ctx context.Context, swipe *Swipe) error
	GetSwipesByUserID(ctx context.Context, userID string) ([]Swipe, error)
	GetSwipes(ctx context.Context, swipe_id string) ([]Swipe, error)
}

type swipeRepository struct {
	db *gorm.DB
}

func NewSwipeRepository(db *gorm.DB) SwipeRepository {
	return &swipeRepository{db: db}
}

func (r *swipeRepository) GetSwipes(ctx context.Context, swipe_id string) ([]Swipe, error) {
	var swipes []Swipe
	err := r.db.WithContext(ctx).
		Where("id = ?", swipe_id).
		Find(&swipes).Error
	return swipes, err
}

func (r *swipeRepository) GetSwipesByUserID(ctx context.Context, userID string) ([]Swipe, error) {
	var swipes []Swipe
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&swipes).Error
	return swipes, err
}

func (r *swipeRepository) Save(ctx context.Context, swipe *Swipe) error {
	swipe.ID = uuid.New().String()
	swipe.CreatedAt = time.Now()

	return r.db.WithContext(ctx).Create(swipe).Error
}
