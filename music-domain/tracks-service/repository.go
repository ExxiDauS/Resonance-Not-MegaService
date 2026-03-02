package tracksservice

import (
	"math"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	GetAllTracks(page int, limit int, name *string, genre *string, artists *string) ([]Track, Pagination, error)
	GetTrackByID(id string) (*Track, error)
	CreateTrack(track *Track) error
	UpdateTrack(track *Track) error
	DeleteTrack(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAllTracks(page int, limit int, name *string, genre *string, artists *string) ([]Track, Pagination, error) {
	var tracks []Track
	var totalCount int64

	query := r.db.Model(&Track{})

	if name != nil && *name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(*name)+"%")
	}

	// If genre is provided, exact match or partial match depending on your needs
	if genre != nil && *genre != "" {
		query = query.Where("LOWER(genre) LIKE ?", "%"+strings.ToLower(*genre)+"%")
	}

	// If artists is provided, search within the artist column
	if artists != nil && *artists != "" {
		query = query.Where("LOWER(artist) LIKE ?", "%"+strings.ToLower(*artists)+"%")
	}

	// 3. Count total records (must be done after filters but before pagination)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, Pagination{}, err
	}

	// 4. Calculate Offset for Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10 // Default limit
	}
	offset := (page - 1) * limit

	// 5. Apply Pagination and Execute Query
	err := query.Limit(limit).Offset(offset).Find(&tracks).Error
	if err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))
	hasNextPage := page < totalPages
	hasPreviousPage := page > 1

	// 6. Create Pagination Metadata
	pagination := Pagination{
		Page:            page,
		Limit:           limit,
		TotalCount:      int(totalCount),
		TotalPage:       totalPages,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}

	return tracks, pagination, nil
}

func (r *repository) GetTrackByID(id string) (*Track, error) {
	var track Track
	if err := r.db.First(&track, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &track, nil
}

func (r *repository) CreateTrack(track *Track) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(track).Error
	})
}

func (r *repository) UpdateTrack(track *Track) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(track).Error
	})
}

func (r *repository) DeleteTrack(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&Track{}, "id = ?", id).Error
	})
}
