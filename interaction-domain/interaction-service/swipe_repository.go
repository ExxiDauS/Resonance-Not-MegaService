package interactionservice

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type SwipeRepository interface {
	Save(ctx context.Context, swipe *Swipe) error
}

type swipeRepository struct {
	db *gorm.DB
}

func NewSwipeRepository(db *gorm.DB) SwipeRepository {
	return &swipeRepository{db: db}
}

func (r *swipeRepository) Save(ctx context.Context, swipe *Swipe) error {
	swipe.ID = uuid.New().String()
	swipe.CreatedAt = time.Now()

	return r.db.WithContext(ctx).Create(swipe).Error
}
