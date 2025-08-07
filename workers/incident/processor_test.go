package incident

import (
	"incident-tracker/config"
	"incident-tracker/models"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetAIClassificationIntegration(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Error loading .env file")
	}
	cfg := config.LoadConfig()
	appCtx := config.NewApplicationContext(cfg)

	incident := &models.Incident{
		Title:           "Server is down",
		Description:     "The main web server is not responding to pings.",
		AffectedService: "Web Server",
	}

	res, err := GetAIClassification(appCtx, incident)
	log.Println(err)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
