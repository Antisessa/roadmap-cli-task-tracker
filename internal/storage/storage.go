package storage

import (
	"bufio"
	"cli-task-tracker/internal/task"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TaskReaderWriter struct {
	Path string
}

func (trw TaskReaderWriter) WriteTask(newTask task.Task) (err error) {
	filename := strconv.Itoa(newTask.Id) + ".json"
	var file *os.File

	// Если файл существует - перезаписываем его, если нет - создаем
	file, err = os.Create(filepath.Join(trw.Path, filename))

	if err != nil {
		return err
	}

	// Изменение именованного результата через defer
	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			err = errors.Join(err, cerr)
		}
	}(file)

	err = trw.writeToFile(newTask, file)
	return err
}

func (trw TaskReaderWriter) ReadTask(id int) (task.Task, error) {
	filename := strconv.Itoa(id) + ".json"
	var res task.Task

	data, err := os.ReadFile(filepath.Join(trw.Path, filename))
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}

	return res, err
}

func (trw TaskReaderWriter) ReadAllByFilter(filterStatus task.Status) ([]task.Task, error) {
	files, err := os.ReadDir(trw.Path)
	if err != nil {
		return []task.Task{}, err
	}

	errFlag := false
	tasks := make([]task.Task, 0, len(files))

	for _, file := range files {
		fileName := file.Name()
		extension := filepath.Ext(fileName)

		switch extension {
		case ".gitkeep":
			continue
		case ".json":
			clearName, found := strings.CutSuffix(fileName, ".json")
			if !found {
				fmt.Println("Cannot clear name from extension")
				_ = deleteFile(filepath.Join(trw.Path, fileName))
				errFlag = true
			}

			var taskId int
			taskId, err = strconv.Atoi(clearName)
			if err != nil {
				fmt.Println("Cannot parse taskId")
				_ = deleteFile(filepath.Join(trw.Path, fileName))
				errFlag = true
				continue
			}

			var newTask task.Task
			newTask, err = trw.ReadTask(taskId)
			if err != nil {
				fmt.Println("Cannot read task from file")
				_ = deleteFile(filepath.Join(trw.Path, fileName))
				errFlag = true
				continue
			}
			if filterStatus == task.None || filterStatus == newTask.Status {
				tasks = append(tasks, newTask)
			}

		default:
			fmt.Println("Unknown file extension")
			_ = deleteFile(filepath.Join(trw.Path, fileName))
			errFlag = true
		}
	}

	if len(tasks) == 0 && errFlag {
		err = errors.New("all tasks was corrupted")
	} else if len(tasks) == 0 && len(files) > 1 {
		fmt.Println("Found no tasks")
	}

	return tasks, err
}

func (trw TaskReaderWriter) writeToFile(task task.Task, file *os.File) error {
	writer := bufio.NewWriter(file)

	bytes, err := task.ToJson()
	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		cerr := deleteFile(filepath.Join(trw.Path, file.Name()))
		if cerr != nil {
			err = errors.Join(err, cerr)
		}
		return err
	}

	return nil
}

func (trw TaskReaderWriter) LastId() (int, error) {
	files, err := os.ReadDir(trw.Path)
	if err != nil {
		return 0, err
	}

	// Перевернуть слайс в одну строчку
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	res := 0

	for _, file := range files {
		fileName := file.Name()
		extension := filepath.Ext(fileName)

		switch extension {
		case ".gitkeep":
			continue
		case ".json":
			clearFileName, _ := strings.CutSuffix(fileName, ".json")
			res, err = strconv.Atoi(clearFileName)
			if err == nil {
				return res, nil
			}
			fallthrough
		default:
			err = deleteFile(filepath.Join(trw.Path, fileName))
			return res, err
		}
	}
	return res, nil
}

func deleteFile(filePath string) error {
	fmt.Println("Удаляю невалидный файл", filePath)
	return os.Remove(filePath)
}
