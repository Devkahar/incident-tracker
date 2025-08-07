package controller

import (
	"incident-tracker/config"
	"incident-tracker/errors"
	"incident-tracker/models"
	"incident-tracker/repository"
	"incident-tracker/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateIncidentRequest struct {
	Title           string `json:"title" binding:"required,min=3,max=20"`
	Description     string `json:"description" binding:"required,min=3,max=1000"`
	AffectedService string `json:"affected_service" binding:"required,min=3,max=50"`
}

func CreateIncident(appContext *config.ApplicationContext, c *gin.Context) (*models.Incident, error) {
	var req CreateIncidentRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return nil, err
	}

	incident := &models.Incident{
		Title:           req.Title,
		Description:     req.Description,
		AffectedService: req.AffectedService,
		RequestStatus:   models.StatusPending,
		AIModel:         models.ModelOpenAI,
	}

	repo := repository.NewIncidentRepository(appContext.DB)
	if err := repo.CreateIncident(incident); err != nil {
		appContext.Logger.Sugar().Errorf("Failed to create incident: %v", err)
		return nil, err
	}

	return incident, nil
}

func GetIncidents(appContext *config.ApplicationContext, c *gin.Context) (*[]models.Incident, error) {
	// Query parameters
	search := c.Query("search")
	requestStatus := c.Query("request_status")
	aiModel := c.Query("ai_model")
	aiSeverity := c.Query("ai_severity")
	aiCategory := c.Query("ai_category")

	// Validate RequestStatus
	if requestStatus != "" {
		validStatuses := map[models.RequestStatus]bool{
			models.StatusPending:    true,
			models.StatusInProgress: true,
			models.StatusProcessed:  true,
			models.StatusFailed:     true,
			models.StatusCompleted:  true,
		}
		if !validStatuses[models.RequestStatus(requestStatus)] {
			return nil, errors.NewAPIError(http.StatusBadRequest, "Invalid request_status", nil)
		}
	}

	// Validate AIModel
	if aiModel != "" {
		validModels := map[models.AIModel]bool{
			models.ModelOpenAI: true,
			models.ModelGroq:   true,
			models.ModelGemini: true,
		}
		if !validModels[models.AIModel(aiModel)] {
			return nil, errors.NewAPIError(http.StatusBadRequest, "Invalid ai_model", nil)
		}
	}

	// Pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		return nil, errors.NewAPIError(http.StatusBadRequest, "Invalid page number", err)
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		return nil, errors.NewAPIError(http.StatusBadRequest, "Invalid limit", err)
	}
	offset := (page - 1) * limit

	// Filters
	filters := &repository.IncidentFilters{
		Search:        search,
		RequestStatus: models.RequestStatus(requestStatus),
		AIModel:       models.AIModel(aiModel),
		AISeverity:    aiSeverity,
		AICategory:    aiCategory,
	}

	// Repository
	repo := repository.NewIncidentRepository(appContext.DB)
	incidents, err := repo.GetIncidents(filters, limit, offset)
	if err != nil {
		appContext.Logger.Sugar().Errorf("Failed to fetch incidents: %v", err)
		return nil, errors.NewAPIError(http.StatusInternalServerError, "Failed to retrieve incidents", err)
	}

	return &incidents, nil
}
