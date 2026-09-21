package intent

import "errors"

var (
	ErrNameCannotBeEmpty = errors.New("name cannot be empty")
	ErrIntentIsNotActive = errors.New("intent is not active")
	ErrIntentIsNotPaused = errors.New("intent is not paused")
)
