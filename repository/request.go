package repository

import (
	"incident-tracker/models"

	"gorm.io/gorm"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) CreateRequest(request *models.Request) error {
	if err := r.db.Create(request).Error; err != nil {
		return err
	}
	return nil
}
