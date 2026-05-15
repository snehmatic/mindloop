package task

import "errors"

const maxTitleLength = 100

var (
	ErrTaskTitleEmpty       = errors.New("task title cannot be empty")
	ErrTaskTitleTooLong     = errors.New("task title exceeds maximum length")
	ErrSubTaskTitleEmpty    = errors.New("subtask title cannot be empty")
	ErrSubTaskTitleTooLong  = errors.New("subtask title exceeds maximum length")
	ErrSubTaskInvalidTaskID = errors.New("subtask must have a valid task ID")
)

func ValidateTaskTitle(title string) error {
	if title == "" {
		return ErrTaskTitleEmpty
	}
	if len(title) > maxTitleLength {
		return ErrTaskTitleTooLong
	}
	return nil
}

func ValidateSubTaskTitle(title string) error {
	if title == "" {
		return ErrSubTaskTitleEmpty
	}
	if len(title) > maxTitleLength {
		return ErrSubTaskTitleTooLong
	}
	return nil
}

func ValidateTaskID(taskID uint) error {
	if taskID == 0 {
		return ErrSubTaskInvalidTaskID
	}
	return nil
}
