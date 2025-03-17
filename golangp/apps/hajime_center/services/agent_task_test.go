package services_test

import (
	"hajime/golangp/apps/hajime_center/models"
	"hajime/golangp/apps/hajime_center/services"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&models.AgentTask{})
	if err != nil {
		panic("failed to migrate database")
	}
	return db
}

func TestAgentTaskService(t *testing.T) {
	// Setup test handlers
	handlers := map[string]func() error{
		"testTask": func() error {
			return nil
		},
		"failingTask": func() error {
			return assert.AnError
		},
	}

	tests := []struct {
		name string
		run  func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB)
	}{
		{
			name: "AddTask_Success",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				execTime := time.Now().Add(time.Second)
				err := svc.AddTask(id, "testTask", execTime)
				assert.NoError(t, err)

				var task models.AgentTask
				assert.NoError(t, db.First(&task, "id = ?", id).Error)
				assert.Equal(t, "Pending", task.State)
				assert.Equal(t, "testTask", task.FunctionName)
			},
		},
		{
			name: "AddTask_DuplicateID",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				execTime := time.Now().Add(time.Second)

				err := svc.AddTask(id, "testTask", execTime)
				assert.NoError(t, err)

				err = svc.AddTask(id, "testTask", execTime)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "task ID already exists")
			},
		},
		{
			name: "SetDailyTask_Success",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				err := svc.SetDailyTask(id, "testTask", 23, 59)
				assert.NoError(t, err)

				var task models.AgentTask
				assert.NoError(t, db.First(&task, "id = ?", id).Error)
				assert.Equal(t, "Pending", task.State)
				assert.Equal(t, "testTask", task.FunctionName)
				assert.Equal(t, "24h0m0s", task.Interval)
			},
		},
		{
			name: "SetRecurringTask_InvalidInterval",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				err := svc.SetRecurringTask(id, "testTask", time.Now(), -time.Second)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "interval must be positive")
			},
		},
		{
			name: "CancelTask_Success",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				err := svc.AddTask(id, "testTask", time.Now().Add(time.Second))
				assert.NoError(t, err)

				err = svc.CancelTask(id)
				assert.NoError(t, err)

				var task models.AgentTask
				assert.NoError(t, db.First(&task, "id = ?", id).Error)
				assert.Equal(t, "Cancelled", task.State)
			},
		},
		{
			name: "ExecuteTask_OneTime",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				err := svc.AddTask(id, "testTask", time.Now())
				assert.NoError(t, err)

				// Wait for task to execute
				time.Sleep(2 * time.Second)

				var task models.AgentTask
				assert.NoError(t, db.First(&task, "id = ?", id).Error)
				assert.Equal(t, "Completed", task.State)
				assert.Empty(t, task.Interval)
				assert.False(t, task.LastExecution.IsZero())
				// t.Errorf("%+v\n", task)
			},
		},
		{
			name: "ExecuteTask_Recurring",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				interval := 2 * time.Second
				err := svc.SetRecurringTask(id, "testTask", time.Now(), interval)
				assert.NoError(t, err)

				// Wait for multiple executions
				time.Sleep(3 * time.Second)

				var task models.AgentTask
				assert.NoError(t, db.First(&task, "id = ?", id).Error)
				assert.Equal(t, "Pending", task.State) // Should be pending for next execution
				assert.Equal(t, interval.String(), task.Interval)
				assert.False(t, task.LastExecution.IsZero())
				t.Logf("%+v\n", task)
				assert.True(t, task.ExecutionTime.After(time.Now()))
			},
		},
		{
			name: "ExecuteTask_Failure",
			run: func(t *testing.T, svc *services.AgentTaskService, db *gorm.DB) {
				id := uuid.New().String()
				err := svc.AddTask(id, "failingTask", time.Now())
				assert.NoError(t, err)

				time.Sleep(2 * time.Second)

				var task models.AgentTask
				assert.NoError(t, db.First(&task, "id = ?", id).Error)
				assert.Equal(t, "Pending", task.State) // Should revert to pending on failure
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB()
			svc := services.NewAgentTaskService(db, handlers)
			defer svc.Stop()
			tt.run(t, svc, db)
		})
	}
}

// TestConcurrency tests concurrent task execution
func TestConcurrency(t *testing.T) {
	db := setupTestDB()
	executionCount := 0
	handlers := map[string]func() error{
		"slowTask": func() error {
			time.Sleep(100 * time.Millisecond)
			executionCount++
			return nil
		},
	}
	svc := services.NewAgentTaskService(db, handlers)
	defer svc.Stop()

	// Add multiple tasks with same ID to test concurrency protection
	id := uuid.New().String()
	for i := 0; i < 5; i++ {
		err := svc.AddTask(id, "slowTask", time.Now())
		if i == 0 {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err) // Only first should succeed
		}
	}

	time.Sleep(1 * time.Second)
	assert.Equal(t, 1, executionCount, "Task should only execute once despite multiple attempts")
}
