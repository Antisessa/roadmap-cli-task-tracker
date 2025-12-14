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

func deleteFile(fileName string) error {
	fmt.Println("Удаляю невалидный файл", fileName)
	return os.Remove(fileName)
}
