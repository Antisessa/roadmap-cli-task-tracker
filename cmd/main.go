package main

import (
	"cli-task-tracker/internal/config"
	"cli-task-tracker/internal/storage"
	"cli-task-tracker/internal/task"
	"flag"
	"fmt"
	"strconv"
	"time"
)

const dataDir = "../data"

func main() {
	cfg := &config.Config{
		Storage: storage.TaskReaderWriter{Path: dataDir},
	}

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Not enough arguments. See 'help' instruction")
	}

	var err error
	switch args[0] {
	case `add`:
		err = add(args[1:], cfg.Storage)
	case `update`:
		err = update(args[1:], cfg.Storage)
	case `mark-in-progress`, `mark-done`:
		err = setStatus(args, cfg.Storage)
	case `list`:

	default:
		fallthrough
	case `help`:
		printHelp()
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return

}

func printHelp() {
	fmt.Println("Available command list:")
	fmt.Println("add - Adding a new task")
	fmt.Print("Example - task-cli add \"Buy groceries\" \n\n")

	fmt.Println("update - Update existing task")
	fmt.Print("Example - task-cli update 1 \"Buy groceries and cook dinner\" \n\n")

	fmt.Println("delete - Delete existing task")
	fmt.Print("Example - task-cli delete 1 \n\n")

	fmt.Println("mark-in-progress - Marking a task as in progress")
	fmt.Print("Example - task-cli mark-in-progress 1 \n\n")

	fmt.Println("mark-done - Marking a task as done")
	fmt.Print("Example - task-cli mark-done 1 \n\n")

	fmt.Println("list - Listing all tasks")
	fmt.Print("Example - task-cli list \n\n")

	fmt.Println("list - Listing tasks by status")
	fmt.Println("Example - task-cli list `status`")
	fmt.Println("list command supporting these statuses - done, todo, in-progress")
}

// task-cli add "Buy groceries"
func add(args []string, trw storage.TaskReaderWriter) error {
	if len(args) != 1 {
		return fmt.Errorf("error! Check that task's name was correct. Tip -> put it in \"quotes\"")
	}

	var err error

	newId, err := trw.LastId()
	if err != nil {
		return err
	}

	newTask := task.Task{
		Id:          newId + 1,
		Description: args[0],
		Status:      task.Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = trw.WriteTask(newTask)
	return err
}

// task-cli update 1 "Buy groceries and cook dinner"
func update(args []string, trw storage.TaskReaderWriter) error {
	if len(args) != 2 {
		return fmt.Errorf("error! Check format 'task_id newTaskName'. Tip -> put task's name in \"quotes\"")
	}

	taskId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	Task, err := trw.ReadTask(taskId)
	if err != nil {
		return err
	}

	Task.UpdateDescription(args[1])
	err = trw.WriteTask(Task)
	if err == nil {
		fmt.Println("Successfully updated!")
	}
	return err
}

// task-cli mark-in-progress 1
// task-cli mark-done 1
func setStatus(args []string, trw storage.TaskReaderWriter) error {
	if len(args) != 2 {
		return fmt.Errorf("error! Check format 'mark-[in-progress | done] task_id'")
	}
	var err error

	taskId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	Task, err := trw.ReadTask(taskId)
	if err != nil {
		return err
	}

	if args[0] == `mark-in-progress` {
		Task.UpdateStatus(task.InProgress)
	} else {
		Task.UpdateStatus(task.Done)
	}

	err = trw.WriteTask(Task)
	return err
}
