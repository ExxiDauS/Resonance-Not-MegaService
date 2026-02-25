package userservice

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

func (r *Repository) CreateUserProfile(profile *Profile) error {
	return r.db.Create(profile).Error
}

func (r *Repository) GetAllUserProfiles() ([]Profile, error) {
	var profiles []Profile
	err := r.db.Find(&profiles).Error
	return profiles, err
}

func (r *Repository) GetUserProfileByID(userID uuid.UUID) (*Profile, error) {
	var profile Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if err == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}
	return &profile, err
}

func (r *Repository) UpdateUserProfile(userID uuid.UUID, updates *UpdateProfileInput) (*Profile, error) {
	result := r.db.Model(&Profile{}).Where("user_id = ?", userID).Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Fetch and return the updated profile
	var profile Profile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *Repository) DeleteUserProfile(userID uuid.UUID) error {
	result := r.db.Where("user_id = ?", userID).Delete(&Profile{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
