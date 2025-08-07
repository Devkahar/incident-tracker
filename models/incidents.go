package models

import (
	"time"
)

// RequestStatus enum
type RequestStatus string

const (
	StatusPending    RequestStatus = "PENDING"
	StatusInProgress RequestStatus = "IN_PROGRESS"
	StatusProcessed  RequestStatus = "PROCESSED"
	StatusFailed     RequestStatus = "FAILED"
	StatusCompleted  RequestStatus = "COMPLETED"
)

// AIModel enum
type AIModel string

const (
	ModelOpenAI AIModel = "OPENAI"
	ModelGroq   AIModel = "GROQ"
	ModelGemini AIModel = "GEMINI"
)

type Incident struct {
	ID              uint          `gorm:"primaryKey"`
	Title           string        `gorm:"type:varchar(20);not null"`
	Description     string        `gorm:"type:varchar(1000);not null"`
	AffectedService string        `gorm:"type:varchar(50);not null"`
	RequestStatus   RequestStatus `gorm:"type:varchar(20);default:'PENDING'"`
	AIModel         AIModel       `gorm:"type:varchar(20);default:'OPENAI'"`
	AISeverity      string        `gorm:"type:varchar(20)"`
	AICategory      string        `gorm:"type:varchar(20)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
