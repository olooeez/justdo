package main

import (
	"log"
	"os"

	"github.com/olooeez/justdo/internal/cli"
)

func main() {
	tasksFilePath := os.Getenv("HOME") + "/.justdo.csv"

	app := cli.TaskManagerCli(tasksFilePath)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
