package note

import "errors"

var (
	ErrNoteMustHaveATitleOrContent = errors.New("note must have a title or content")
)
