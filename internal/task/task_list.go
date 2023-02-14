package task

import (
	"encoding/csv"
	"os"

	"github.com/google/uuid"
)

const (
	csvId int = iota
	csvTitle
	csvDescription
	csvStatus
)

type TaskList struct {
	Path  string
	Tasks []*Task
}

func NewTaskList(path string) (*TaskList, error) {
	taskList := &TaskList{Path: path}

	if err := taskList.loadTasks(); err != nil {
		return nil, err
	}

	return taskList, nil
}

func (taskList *TaskList) loadTasks() error {
	file, err := os.OpenFile(taskList.Path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range records {
		id, err := uuid.Parse(row[csvId])
		if err != nil {
			return err
		}

		taskList.Tasks = append(taskList.Tasks, &Task{
			Id:          id,
			Title:       row[csvTitle],
			Description: row[csvDescription],
			Status:      Progress(row[csvStatus]),
		})
	}

	return nil
}

func (taskList *TaskList) saveTasks() error {
	file, err := os.OpenFile(taskList.Path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	for _, t := range taskList.Tasks {
		if err = writer.Write([]string{
			t.Id.String(),
			t.Title,
			t.Description,
			string(t.Status),
		}); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

func (taskList *TaskList) Add(task *Task) error {
	taskList.Tasks = append(taskList.Tasks, task)
	return taskList.saveTasks()
}

func (taskList *TaskList) Remove(task *Task) error {
	for i, currentTask := range taskList.Tasks {
		if task.Id == currentTask.Id {
			taskList.Tasks = append(taskList.Tasks[:i], taskList.Tasks[i+1:]...)
			break
		}
	}

	return taskList.saveTasks()
}

func (taskList *TaskList) Update(oldTask, newTask *Task) error {
	for i, currentTask := range taskList.Tasks {
		if oldTask.Id == currentTask.Id {
			taskList.Tasks[i] = newTask
			break
		}
	}

	return taskList.saveTasks()
}

func (taskList *TaskList) Find(id uuid.UUID) *Task {
	for _, currentTask := range taskList.Tasks {
		if currentTask.Id == id {
			return currentTask
		}
	}

	return nil
}
