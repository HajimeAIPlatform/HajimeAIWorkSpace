package services

import (
	"errors"
	"hajime/golangp/apps/hajime_center/models"

	"gorm.io/gorm"
)

// AgentLogService handles logging operations
type AgentLogService struct {
	db *gorm.DB
}

// NewAgentLogService creates a new instance of AgentLogService
func NewAgentLogService(db *gorm.DB) *AgentLogService {
	return &AgentLogService{db: db}
}

// AddNormalLog adds a normal type log
func (s *AgentLogService) AddNormalLog(message string) error {
	return s.addLog("Normal", message)
}

// AddErrorLog adds an error type log
func (s *AgentLogService) AddErrorLog(message string) error {
	return s.addLog("Error", message)
}

// AddReportingLog adds a reporting type log
func (s *AgentLogService) AddReportingLog(message string) error {
	return s.addLog("Reporting", message)
}

// addLog private helper method to add logs
func (s *AgentLogService) addLog(logType, message string) error {
	if message == "" {
		return errors.New("message cannot be empty")
	}

	log := models.AgentLog{
		Type:    logType,
		Message: message,
	}

	return s.db.Create(&log).Error
}

// GetLogs retrieves logs by type with pagination
func (s *AgentLogService) GetLogs(logType string, page, pageSize int) ([]models.AgentLog, int64, error) {
	var logs []models.AgentLog
	var total int64

	query := s.db.Model(&models.AgentLog{})

	if logType != "" {
		query = query.Where("type = ?", logType)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// Fetch logs
	if err := query.Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
