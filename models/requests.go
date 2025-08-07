package models

import (
	"time"
)

type Request struct {
	ID             uint   `gorm:"primaryKey"`
	IncidentID     uint   `gorm:"not null;index"`     // foreign key index
	RequestBody    string `gorm:"type:text;not null"` // store JSON or prompt
	ResponseBody   string `gorm:"type:text"`          // store JSON from AI
	ResponseStatus uint   `gorm:"type:int"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	Incident Incident `gorm:"foreignKey:IncidentID;constraint:OnDelete:CASCADE"`
}
