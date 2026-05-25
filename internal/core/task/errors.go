package task

import "errors"

var (
	ErrorSubTaskNotFound = errors.New("subtask not found")
	ErrorTaskNotFound    = errors.New("task not found")
)
