package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Status int

const (
	todo = iota
	inProgress
	done
)

type task struct {
	Id          int
	Description string
	Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

/*
task-cli add "Buy groceries"

task-cli update 1 "Buy groceries and cook dinner"
*/
func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Not enough arguments. See 'help' instruction")
	}

	switch args[0] {
	case `help`:
		printHelp()
	case `add`:
		handleAdd(args[1:])
	case `update`:

	case `mark-in-progress`:

	case `mark-done`:

	case `list`:
	}
}

func deleteFile(fileName string) {
	fmt.Println("Удаляю невалидный файл", fileName)
	err := os.Remove(fileName)
	if err != nil {
		panic(err.Error())
	}
}

func getLastId() int {
	const dataDir = "../data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		panic("Cannot open data's folder")
	}

	slices.Reverse(files)

	for _, file := range files {
		fileName := file.Name()
		extension := filepath.Ext(fileName)

		switch extension {
		case ".gitkeep":
			continue
		case ".json":
			clearFileName, _ := strings.CutSuffix(fileName, ".json")
			res, err := strconv.Atoi(clearFileName)
			if err != nil {
				deleteFile(filepath.Join(dataDir, fileName))
			} else {
				return res
			}
		default:
			deleteFile(filepath.Join(dataDir, fileName))
		}
	}
	panic("Found no file, check data folder")
}

func handleAdd(args []string) {
	if len(args) > 1 {
		fmt.Println("Error! Check that task's name was correct. Tip -> put it in \"quotes\"")
	}
	fmt.Println(getLastId())

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
