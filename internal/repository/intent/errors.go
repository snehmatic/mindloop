package intent

import "errors"

var ErrNameCannotBeEmpty = errors.New("name cannot be empty")
var ErrIntentNotActive = errors.New("intent is not active")
var ErrIntentNotPaused = errors.New("intent is not paused")
