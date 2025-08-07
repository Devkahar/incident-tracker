package config

import (
	"incident-tracker/models"
	"log"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ApplicationContext struct {
	DB           *gorm.DB
	Logger       *zap.Logger
	Config       *Config
	OpenAIAPIKey string
	HttpClient   *http.Client
}

func NewApplicationContext(config *Config) *ApplicationContext {
	db, err := gorm.Open(mysql.Open(config.DB.GetDBURL()))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Incident{}, &models.Request{}); err != nil {
		log.Fatal(err)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	return &ApplicationContext{
		DB:           db,
		Config:       config,
		Logger:       logger,
		OpenAIAPIKey: config.OpenAI.APIKey,
		HttpClient:   &http.Client{},
	}
}
