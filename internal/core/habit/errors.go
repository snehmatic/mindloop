package habit

import "errors"

var (
	ErrHabitCannotBeNil                  = errors.New("habit cannot be nil")
	ErrHabitNotFound                     = errors.New("habit not found")
	ErrHabitAlreadyCompletedForInterval  = errors.New("habit already completed for interval")
	ErrNoExistingLogFoundForThisInterval = errors.New("no existing log found for this interval")
	ErrHabitIsAlreadyMarkedAsUndone      = errors.New("habit is already marked as undone")
)
