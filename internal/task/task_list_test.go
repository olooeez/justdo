package task

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskList(t *testing.T) {
	path := "task_list_test.csv"

	defer os.Remove(path)

	taskList, err := NewTaskList(path)

	assert.NoError(t, err)
	assert.NotNil(t, taskList)
	assert.Equal(t, path, taskList.Path)

	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestLoadTasks(t *testing.T) {
	path := "task_list_test.csv"

	defer os.Remove(path)

	data := []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c8,Test Task 1,This is a test task,Not Started\n" + "6ba7b811-9dad-11d1-80b4-00c04fd430c8,Test Task 2,This is another test task,In Progress\n")

	taskList, err := NewTaskList(path)
	assert.NoError(t, err)

	err = os.WriteFile(path, data, 0666)
	assert.NoError(t, err)

	err = taskList.loadTasks()

	assert.NoError(t, err)
	assert.Len(t, taskList.Tasks, 2)

	assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", taskList.Tasks[0].Id.String())
	assert.Equal(t, "Test Task 1", taskList.Tasks[0].Title)
	assert.Equal(t, "This is a test task", taskList.Tasks[0].Description)
	assert.Equal(t, NotStarted, taskList.Tasks[0].Status)

	assert.Equal(t, "6ba7b811-9dad-11d1-80b4-00c04fd430c8", taskList.Tasks[1].Id.String())
	assert.Equal(t, "Test Task 2", taskList.Tasks[1].Title)
	assert.Equal(t, "This is another test task", taskList.Tasks[1].Description)
	assert.Equal(t, InProgress, taskList.Tasks[1].Status)
}

func TestSaveTasks(t *testing.T) {
	path := "task_list_test.csv"
	defer os.Remove(path)

	tl, err := NewTaskList(path)
	assert.NoError(t, err)

	task1 := NewTask("Task 1", "Description 1")
	task2 := NewTask("Task 1", "Description 1")

	tl.Tasks = append(tl.Tasks, task1)
	tl.Tasks = append(tl.Tasks, task2)

	err = tl.saveTasks()
	assert.NoError(t, err)

	tl, err = NewTaskList(path)
	assert.NoError(t, err)
	assert.Len(t, tl.Tasks, 2)

	assert.Equal(t, tl.Tasks[0].Id, task1.Id)
	assert.Equal(t, tl.Tasks[0].Title, task1.Title)
	assert.Equal(t, tl.Tasks[0].Description, task1.Description)
	assert.Equal(t, tl.Tasks[0].Status, task1.Status)

	assert.Equal(t, tl.Tasks[1].Id, task2.Id)
	assert.Equal(t, tl.Tasks[1].Title, task2.Title)
	assert.Equal(t, tl.Tasks[1].Description, task2.Description)
	assert.Equal(t, tl.Tasks[1].Status, task2.Status)
}

func TestAdd(t *testing.T) {
	path := "task_list_test.csv"

	tl, err := NewTaskList(path)
	assert.NoError(t, err)

	defer os.Remove(path)

	task1 := NewTask("Test Task 1", "This is a test task")
	task2 := NewTask("Test Task 2", "This is another test task")

	err = tl.Add(task1)
	assert.NoError(t, err)

	err = tl.Add(task2)
	assert.NoError(t, err)

	assert.Equal(t, len(tl.Tasks), 2)
}

func TestRemove(t *testing.T) {
	path := "task_list_test.csv"

	tl, err := NewTaskList(path)
	assert.NoError(t, err)

	defer os.Remove(path)

	task := NewTask("Test Task", "This is a test task")

	err = tl.Add(task)
	assert.NoError(t, err)

	assert.Equal(t, len(tl.Tasks), 1)

	err = tl.Remove(task)
	assert.NoError(t, err)

	assert.Equal(t, len(tl.Tasks), 0)
}

func TestUpdate(t *testing.T) {
	path := "task_list_test.csv"

	tl, err := NewTaskList(path)
	assert.NoError(t, err)

	defer os.Remove(path)

	oldTask := NewTask("Test Task", "This is a test task")

	err = tl.Add(oldTask)
	assert.NoError(t, err)

	assert.Equal(t, len(tl.Tasks), 1)
	assert.Equal(t, tl.Tasks[0].Title, "Test Task")

	newTask := NewTask("Updated Test Task", "Updated this is a test task")
	newTask.Id = oldTask.Id

	err = tl.Update(oldTask, newTask)
	assert.NoError(t, err)

	assert.Equal(t, len(tl.Tasks), 1)
	assert.Equal(t, tl.Tasks[0].Title, "Updated Test Task")
	assert.Equal(t, tl.Tasks[0].Description, "Updated this is a test task")
}

func TestFind(t *testing.T) {
	path := "task_list_test.csv"

	tl, err := NewTaskList(path)
	assert.NoError(t, err)

	defer os.Remove(path)

	task1 := NewTask("Test Task 1", "This is a test task")
	taskNotStored := NewTask("Test Task 2", "This is another test task")

	err = tl.Add(task1)
	assert.NoError(t, err)

	taskFound := tl.Find(task1.Id)
	assert.NotNil(t, taskFound)

	taskNotFound := tl.Find(taskNotStored.Id)
	assert.Nil(t, taskNotFound)
}
