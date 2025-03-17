package models

import (
	"time"

	"gorm.io/gorm"
)

type AgentTask struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	State         string         `gorm:"type:varchar(20);not null" json:"state"`
	FunctionName  string         `gorm:"type:varchar(100);not null" json:"function_name"`
	ExecutionTime time.Time      `gorm:"not null" json:"execution_time"`
	Interval      string         `gorm:"type:varchar(20)" json:"interval,omitempty"`
	LastExecution time.Time      `json:"last_execution,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
