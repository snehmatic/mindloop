package journal

import "errors"

var ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
var ErrContentCannotBeEmpty = errors.New("content cannot be empty")
