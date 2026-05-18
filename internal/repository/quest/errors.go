package quest

import "errors"

var ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
var ErrQuestAlreadyActive = errors.New("a side quest is already active")
var ErrQuestNotActive = errors.New("side quest is not active")
