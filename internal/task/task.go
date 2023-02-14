package task

import "github.com/google/uuid"

type Progress string

const (
	NotStarted Progress = "Not Started"
	InProgress Progress = "In Progress"
	Completed  Progress = "Completed"
)

type Task struct {
	Id          uuid.UUID
	Title       string
	Description string
	Status      Progress
}

func NewTask(title, description string) *Task {
	return &Task{
		Id:          uuid.New(),
		Title:       title,
		Description: description,
		Status:      NotStarted,
	}
}
