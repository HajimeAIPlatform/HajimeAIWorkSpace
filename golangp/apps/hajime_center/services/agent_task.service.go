package services

import (
	"context"
	"errors"
	"hajime/golangp/apps/hajime_center/models"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// AgentTaskService handles task operations
type AgentTaskService struct {
	db           *gorm.DB
	taskHandlers map[string]func() error
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	executing    sync.Map // To track currently executing tasks
	dbMutex      sync.Mutex
}

// NewAgentTaskService creates a new instance
func NewAgentTaskService(db *gorm.DB, handlers map[string]func() error) *AgentTaskService {
	ctx, cancel := context.WithCancel(context.Background())
	svc := &AgentTaskService{
		db:           db,
		taskHandlers: handlers,
		ctx:          ctx,
		cancel:       cancel,
	}
	svc.StartTaskChecker()
	return svc
}

// AddTask adds a one-time task
func (s *AgentTaskService) AddTask(id, functionName string, executionTime time.Time) error {
	return s.addTask(id, functionName, executionTime, "")
}

// SetRecurringTask sets a task with arbitrary recurrence interval
func (s *AgentTaskService) SetRecurringTask(id, functionName string, firstExecution time.Time, interval time.Duration) error {
	if interval <= 0 {
		return errors.New("interval must be positive")
	}

	now := time.Now()
	execTime := firstExecution
	// If first execution is in the past, calculate next occurrence
	for execTime.Before(now) {
		execTime = execTime.Add(interval)
	}

	return s.addTask(id, functionName, execTime, interval.String())
}

// SetDailyTask sets a daily recurring task (convenience method)
func (s *AgentTaskService) SetDailyTask(id, functionName string, hour, minute int) error {
	now := time.Now()
	firstExec := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	return s.SetRecurringTask(id, functionName, firstExec, 24*time.Hour)
}

// SetHourlyTask sets an hourly recurring task (convenience method)
func (s *AgentTaskService) SetHourlyTask(id, functionName string, minute int) error {
	now := time.Now()
	firstExec := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, 0, 0, now.Location())
	return s.SetRecurringTask(id, functionName, firstExec, time.Hour)
}

// addTask internal helper method
func (s *AgentTaskService) addTask(id, functionName string, executionTime time.Time, scheduleType string) error {
	if id == "" || functionName == "" {
		return errors.New("id and function_name cannot be empty")
	}
	if _, exists := s.taskHandlers[functionName]; !exists {
		return errors.New("unknown function name")
	}

	// Check for duplicate ID
	var count int64
	if err := s.db.Model(&models.AgentTask{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("task ID already exists")
	}

	task := models.AgentTask{
		ID:            id,
		State:         "Pending",
		FunctionName:  functionName,
		ExecutionTime: executionTime,
		Interval:      scheduleType,
	}
	return s.db.Create(&task).Error
}

// CancelTask cancels a pending task
func (s *AgentTaskService) CancelTask(id string) error {
	result := s.db.Model(&models.AgentTask{}).
		Where("id = ? AND state = ?", id, "Pending").
		Update("state", "Cancelled")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("task not found or not in pending state")
	}
	return nil
}

// StartTaskChecker starts the task monitoring goroutine
func (s *AgentTaskService) StartTaskChecker() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.ctx.Done():
				return
			case <-ticker.C:
				s.checkAndExecuteTasks()
			}
		}
	}()
}

// Stop stops the task checker
func (s *AgentTaskService) Stop() {
	s.cancel()
	s.wg.Wait()
}

// checkAndExecuteTasks checks for due tasks
func (s *AgentTaskService) checkAndExecuteTasks() {
	var tasks []models.AgentTask
	now := time.Now()

	if err := s.db.Where("state = ? AND execution_time <= ?", "Pending", now).
		Find(&tasks).Error; err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return
	}

	for _, task := range tasks {
		s.wg.Add(1)
		go func(t models.AgentTask) {
			defer s.wg.Done()
			if err := s.executeTask(&t); err != nil {
				log.Printf("Error executing task %s: %v", t.ID, err)
			}
		}(task)
	}
}

// executeTask executes a single task with competition prevention
func (s *AgentTaskService) executeTask(task *models.AgentTask) error {
	// Check if task is already executing
	if _, loaded := s.executing.LoadOrStore(task.ID, true); loaded {
		return nil // Task is already being executed
	}
	defer s.executing.Delete(task.ID)

	s.dbMutex.Lock()
	defer s.dbMutex.Unlock()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Double-check state in transaction
	var currentTask models.AgentTask
	if err := tx.Where("id = ? AND state = ?", task.ID, "Pending").First(&currentTask).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update state to Running
	if err := tx.Model(task).Update("state", "Running").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Execute the task
	if handler, ok := s.taskHandlers[task.FunctionName]; ok {
		if err := handler(); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Handle scheduling and completion
	updates := map[string]interface{}{
		"state":          "Completed",
		"last_execution": time.Now(),
	}

	// Handle recurring tasks
	if task.Interval != "" {
		if interval, err := time.ParseDuration(task.Interval); err == nil {
			updates["execution_time"] = task.ExecutionTime.Add(interval)
			updates["state"] = "Pending"
		} else {
			log.Printf("Failed to parse interval %s for task %s: %v", task.Interval, task.ID, err)
		}
	}

	if err := tx.Model(task).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
