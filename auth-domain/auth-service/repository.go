package authservice

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateUser(user *Credential) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetByEmail(email string) (*Credential, error) {
	var user Credential
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *Repository) UpdateLastLogin(user *Credential) error {
	return r.db.Model(user).Update("last_login", gorm.Expr("NOW()")).Error
}