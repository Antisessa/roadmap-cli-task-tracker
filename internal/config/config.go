package config

import "cli-task-tracker/internal/storage"

type Config struct {
	Storage storage.TaskReaderWriter
}
