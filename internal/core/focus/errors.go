package focus

import "errors"

var (
	ErrTitleCannotBeEmpty           = errors.New("title cannot be empty")
	ErrAFocusSessionIsAlreadyActive = errors.New("a focus session is already active")
	ErrFocusSessionIsNotActive      = errors.New("focus session is not active")
	ErrRatingMustBeBetween0And10    = errors.New("rating must be between 0 and 10")
	ErrFocusSessionIsNotEnded       = errors.New("focus session is not ended")
	ErrFocusSessionIsNotPaused      = errors.New("focus session is not paused")
)
