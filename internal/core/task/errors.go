package task

import "errors"

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrSubtaskNotFound = errors.New("subtask not found")
)
