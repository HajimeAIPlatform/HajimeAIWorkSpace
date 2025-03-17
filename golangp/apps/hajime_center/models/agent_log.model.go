package models

import (
	"gorm.io/gorm"
)

// AgentLog represents the logging table structure
type AgentLog struct {
	gorm.Model
	Type    string `gorm:"type:varchar(20);not null" json:"type"` // Normal, Error, Reporting
	Message string `gorm:"type:text;not null" json:"message"`
}
