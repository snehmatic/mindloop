package journal

import "errors"

var (
	ErrTitleCannotBeEmpty   = errors.New("title cannot be empty")
	ErrContentCannotBeEmpty = errors.New("content cannot be empty")
)
