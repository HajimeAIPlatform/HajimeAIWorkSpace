package test

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"hajime/golangp/apps/AgentTown/task"
	"testing"
)

func TestSmokeTest(t *testing.T) {
	assert.Equal(t, 10, 9+1)
}

func TestSmokeTest2(t *testing.T) {
	assert.Equal(t, 12, 9+3)
}

func TestTaskGeneration(t *testing.T) {
	tsk := task.NewTask("Test Task")
	fmt.Printf("Task ID: %s\n\n", tsk.ID)
	assert.Equal(t, "Test Task", tsk.Description)
}
