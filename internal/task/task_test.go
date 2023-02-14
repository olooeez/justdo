package task

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	task := NewTask("Test Task", "This is a test task")

	assert.NotNil(t, task)
	assert.IsType(t, &Task{}, task)
	assert.NotEqual(t, uuid.Nil, task.Id)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "This is a test task", task.Description)
	assert.Equal(t, NotStarted, task.Status)
}
