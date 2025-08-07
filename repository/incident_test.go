package repository

import (
	"incident-tracker/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&models.Incident{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func TestIncidentRepository(t *testing.T) {
	db := setup(t)
	repo := NewIncidentRepository(db)

	t.Run("CreateIncident", func(t *testing.T) {
		incident := &models.Incident{
			Title:           "Test Incident",
			Description:     "Test Description",
			AffectedService: "Test Service",
		}
		err := repo.CreateIncident(incident)
		assert.NoError(t, err)
		assert.NotZero(t, incident.ID)
	})

	t.Run("GetIncidents", func(t *testing.T) {
		// Clear the table
		db.Exec("DELETE FROM incidents")

		// Create some incidents
		repo.CreateIncident(&models.Incident{Title: "Incident 1"})
		repo.CreateIncident(&models.Incident{Title: "Incident 2"})

		incidents, err := repo.GetIncidents(&IncidentFilters{}, 1, 0)
		assert.NoError(t, err)
		assert.Len(t, incidents, 1)

		incidents, err = repo.GetIncidents(&IncidentFilters{}, 0, 1)
		assert.NoError(t, err)
		assert.Len(t, incidents, 1)
	})

	t.Run("GetIncidentByID", func(t *testing.T) {
		incident := &models.Incident{Title: "Test Incident"}
		repo.CreateIncident(incident)

		found, err := repo.GetIncidentByID(incident.ID)
		assert.NoError(t, err)
		assert.Equal(t, incident.Title, found.Title)
	})

	t.Run("GetIncidentByRequestStatusAndAIModel", func(t *testing.T) {
		// Clear the table
		db.Exec("DELETE FROM incidents")
		repo.CreateIncident(&models.Incident{Title: "Incident 1", RequestStatus: "PENDING", AIModel: "OPENAI"})
		repo.CreateIncident(&models.Incident{Title: "Incident 2", RequestStatus: "PROCESSED", AIModel: "OPENAI"})
		repo.CreateIncident(&models.Incident{Title: "Incident 3", RequestStatus: "PENDING", AIModel: "GROQ"})

		incidents, err := repo.GetIncidentByRequestStatusAndAIModel("PENDING", "", 0, 0)
		assert.NoError(t, err)
		assert.Len(t, incidents, 2)

		incidents, err = repo.GetIncidentByRequestStatusAndAIModel("", "OPENAI", 0, 0)
		assert.NoError(t, err)
		assert.Len(t, incidents, 2)

		incidents, err = repo.GetIncidentByRequestStatusAndAIModel("PENDING", "OPENAI", 0, 0)
		assert.NoError(t, err)
		assert.Len(t, incidents, 1)
	})
}
