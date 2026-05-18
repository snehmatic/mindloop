package note

import "errors"

var ErrNoteMustHaveTitleOrContent = errors.New("note must have a title or content")
