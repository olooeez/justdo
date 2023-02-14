package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/google/uuid"
	"github.com/olooeez/justdo/internal/task"
	"github.com/urfave/cli/v2"
)

var (
	quiet               bool
	globalTasksFilePath string
)

func TaskManagerCli(tasksFilePath string) *cli.App {
	app := cli.NewApp()

	globalTasksFilePath = tasksFilePath

	app.HideHelpCommand = true

	app.Name = "justdo"
	app.Usage = "simple cli todo list application"
	app.Version = "0.1.0"

	app.Authors = []*cli.Author{
		{
			Name:  "Luiz Felipe",
			Email: "luizfelipecastrovb@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "quiet",
			Usage:       "enable quiet mode",
			Destination: &quiet,
			Value:       false,
			Aliases:     []string{"q"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:      "list",
			Usage:     "list all tasks",
			UsageText: "justdo list",
			Action:    listCommand,
		},
		{
			Name:      "add",
			Usage:     "add a task",
			ArgsUsage: "TITLE DESCRIPTION",
			UsageText: "justdo add TITLE DESCRIPTION",
			Action:    addCommand,
		},
		{
			Name:      "remove",
			Usage:     "remove a task",
			ArgsUsage: "ID",
			UsageText: "justdo remove ID",
			Action:    removeCommand,
		},
		{
			Name:      "update",
			Usage:     "change the status of a task",
			ArgsUsage: "ID STATUS",
			UsageText: "justdo update ID STATUS",
			Action:    updateCommand,
		},
	}

	return app
}

func listCommand(c *cli.Context) error {
	actualNumArgs := c.Args().Len()

	if actualNumArgs != 0 {
		return fmt.Errorf("invalid number of arguments (expected 0, got %d). see 'justdo list --help' for more info", actualNumArgs)
	}

	taskList, err := task.NewTaskList(globalTasksFilePath)
	if err != nil {
		return err
	}

	tab := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.DiscardEmptyColumns)

	fmt.Fprintf(tab, "ID\tTITLE\tDESCRIPTION\tSTATUS\t\n")
	for _, t := range taskList.Tasks {
		fmt.Fprintf(tab, "%s\t%s\t%s\t%s\t\n", t.Id.String(), t.Title, t.Description, t.Status)
	}

	return tab.Flush()
}

func addCommand(c *cli.Context) error {
	actualNumArgs := c.Args().Len()

	if actualNumArgs != 2 {
		return fmt.Errorf("invalid number of arguments (expected 2, got %d). see 'justdo add --help' for more info", actualNumArgs)
	}

	title := c.Args().First()
	if _, err := strconv.Atoi(title); err == nil {
		return errors.New("expected TITLE to be a string, but got a int. see 'justdo add --help' for more info")
	}

	description := c.Args().Get(1)
	if _, err := strconv.Atoi(description); err == nil {
		return errors.New("expected DESCRIPTION to be a string, but got a int. see 'justdo add --help' for more info")
	}

	taskList, err := task.NewTaskList(globalTasksFilePath)
	if err != nil {
		return err
	}

	newTask := task.NewTask(title, description)

	err = taskList.Add(newTask)
	if err != nil {
		return err
	}

	if !quiet {
		log.Printf("added a new task with the title of '%s' and description of '%s'\n", title, description)
	}

	return nil
}

func removeCommand(c *cli.Context) error {
	actualNumArgs := c.Args().Len()

	if actualNumArgs != 1 {
		return fmt.Errorf("invalid number of arguments (expected 1, got %d). see 'justdo remove --help' for more info", actualNumArgs)
	}

	id := c.Args().First()
	if _, err := strconv.Atoi(id); err == nil {
		return errors.New("expected ID to be a string, but got a int. see 'justdo remove --help' for more info")
	}

	taskList, err := task.NewTaskList(globalTasksFilePath)
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	taskToRemove := taskList.Find(uuid)
	if taskToRemove == nil {
		return fmt.Errorf("couldn't find the task with the id of '%s'", id)
	}

	err = taskList.Remove(taskToRemove)
	if err != nil {
		return err
	}

	if !quiet {
		log.Printf("removed a task with the id of '%s'\n", id)
	}

	return nil
}

func updateCommand(c *cli.Context) error {
	actualNumArgs := c.Args().Len()

	if actualNumArgs != 2 {
		return fmt.Errorf("invalid number of arguments (expected 2, got %d). see 'justdo update --help' for more info", actualNumArgs)
	}

	id := c.Args().First()
	if _, err := strconv.Atoi(id); err == nil {
		return errors.New("expected ID to be a string, but got a int. see 'justdo remove --help' for more info")
	}

	status := c.Args().Get(1)
	if status != string(task.NotStarted) && status != string(task.InProgress) && status != string(task.Completed) {
		return fmt.Errorf("'%v' progress is not recognized", status)
	}

	taskList, err := task.NewTaskList(globalTasksFilePath)
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	targetTask := taskList.Find(uuid)
	if targetTask == nil {
		return fmt.Errorf("couldn't find the task with the id of '%s'", id)
	}

	newTask := targetTask
	newTask.Status = task.Progress(status)

	err = taskList.Update(targetTask, newTask)
	if err != nil {
		return err
	}

	if !quiet {
		log.Printf("updated a task with the id of '%s' and status from '%s' to '%s'\n", id, status, newTask.Status)
	}

	return nil
}
