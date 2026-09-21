package quest

import "errors"

var (
	ErrTitleCannotBeEmpty        = errors.New("title cannot be empty")
	ErrASideQuestIsAlreadyActive = errors.New("a side quest is already active")
	ErrSideQuestIsNotActive      = errors.New("side quest is not active")
)
