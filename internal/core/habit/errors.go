package habit

import "errors"

var ErrHabitCannotBeNil = errors.New("habit cannot be nil")
var ErrHabitNotFound = errors.New("habit not found")
var ErrHabitAlreadyCompleted = errors.New("habit already completed for interval")
var ErrNoExistingLog = errors.New("no existing log found for this interval")
var ErrHabitAlreadyUndone = errors.New("habit is already marked as undone")
