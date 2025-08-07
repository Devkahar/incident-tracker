package controller

import (
	"bytes"
	"encoding/json"
	"incident-tracker/config"
	"incident-tracker/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup(t *testing.T) *config.ApplicationContext {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&models.Incident{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return &config.ApplicationContext{
		DB: db,
	}
}

func TestCreateIncident(t *testing.T) {
	appContext := setup(t)
	gin.SetMode(gin.TestMode)

	t.Run("ValidRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := CreateIncidentRequest{
			Title:           "Test Incident",
			Description:     "Test Description",
			AffectedService: "Test Service",
		}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		incident, err := CreateIncident(appContext, c)

		assert.NoError(t, err)
		assert.NotNil(t, incident)
		assert.Equal(t, "Test Incident", incident.Title)
		assert.Equal(t, models.StatusPending, incident.RequestStatus)
		assert.Equal(t, models.ModelOpenAI, incident.AIModel)
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := CreateIncidentRequest{
			Title: "a",
		}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		_, err := CreateIncident(appContext, c)

		assert.Error(t, err)
	})
}

func TestGetIncidents(t *testing.T) {
	appContext := setup(t)
	gin.SetMode(gin.TestMode)

	// Seed one incident for valid case
	appContext.DB.Create(&models.Incident{
		Title:           "Network Failure",
		Description:     "Router not responding",
		AffectedService: "Network",
		RequestStatus:   models.StatusPending,
		AIModel:         models.ModelOpenAI,
		AISeverity:      "High",
		AICategory:      "Network",
	})

	t.Run("ValidRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/?search=network&request_status=PENDING&ai_model=OPENAI&ai_severity=High&ai_category=Network&page=1&limit=10", nil)
		c.Request = req

		incidents, err := GetIncidents(appContext, c)

		assert.NoError(t, err)
		assert.NotNil(t, incidents)
		assert.GreaterOrEqual(t, len(*incidents), 1)
	})

	t.Run("InvalidRequestStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/?request_status=INVALID", nil)
		c.Request = req

		_, err := GetIncidents(appContext, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid request_status")
	})

	t.Run("InvalidAIModel", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/?ai_model=INVALID", nil)
		c.Request = req

		_, err := GetIncidents(appContext, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid ai_model")
	})

	t.Run("InvalidPage", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/?page=zero", nil)
		c.Request = req

		_, err := GetIncidents(appContext, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid page number")
	})

	t.Run("InvalidLimit", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req, _ := http.NewRequest(http.MethodGet, "/?limit=zero", nil)
		c.Request = req

		_, err := GetIncidents(appContext, c)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid limit")
	})
}
