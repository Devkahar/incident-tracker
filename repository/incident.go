package repository

import (
	"incident-tracker/models"

	"gorm.io/gorm"
)

type IncidentRepository struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) *IncidentRepository {
	return &IncidentRepository{db: db}
}

func (r *IncidentRepository) CreateIncident(incident *models.Incident) error {
	if err := r.db.Create(incident).Error; err != nil {
		return err
	}
	return nil
}

type IncidentFilters struct {
	Search        string
	RequestStatus models.RequestStatus
	AIModel       models.AIModel
	AISeverity    string
	AICategory    string
}

func (r *IncidentRepository) GetIncidents(filters *IncidentFilters, limit, offset int) ([]models.Incident, error) {
	var incidents []models.Incident
	query := r.db.Model(&models.Incident{})

	if filters.Search != "" {
		searchQuery := "%" + filters.Search + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR affected_service LIKE ?", searchQuery, searchQuery, searchQuery)
	}
	if filters.RequestStatus != "" {
		query = query.Where("request_status = ?", filters.RequestStatus)
	}
	if filters.AIModel != "" {
		query = query.Where("ai_model = ?", filters.AIModel)
	}
	if filters.AISeverity != "" {
		query = query.Where("ai_severity = ?", filters.AISeverity)
	}
	if filters.AICategory != "" {
		query = query.Where("ai_category = ?", filters.AICategory)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&incidents).Error; err != nil {
		return nil, err
	}
	return incidents, nil
}

func (r *IncidentRepository) GetIncidentByID(id uint) (*models.Incident, error) {
	var incident models.Incident
	if err := r.db.First(&incident, id).Error; err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *IncidentRepository) GetIncidentByRequestStatusAndAIModel(status models.RequestStatus, model models.AIModel, limit, offset int) ([]models.Incident, error) {
	var incidents []models.Incident
	query := r.db.Model(&models.Incident{})

	if status != "" {
		query = query.Where("request_status = ?", status)
	}

	if model != "" {
		query = query.Where("ai_model = ?", model)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&incidents).Error; err != nil {
		return nil, err
	}

	return incidents, nil
}

func (r *IncidentRepository) UpdateIncident(incident *models.Incident) error {
	if err := r.db.Save(incident).Error; err != nil {
		return err
	}
	return nil
}
