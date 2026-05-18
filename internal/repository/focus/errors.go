package focus

import "errors"

var ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
var ErrSessionAlreadyActive = errors.New("a focus session is already active")
var ErrSessionNotActive = errors.New("focus session is not active")
var ErrRatingOutOfRange = errors.New("rating must be between 0 and 10")
var ErrSessionNotEnded = errors.New("focus session is not ended")
var ErrSessionNotPaused = errors.New("focus session is not paused")
