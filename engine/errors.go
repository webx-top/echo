package engine

import "errors"

var (
	ErrUnsupported      = errors.New(`Unsupported`)
	ErrAlreadyCommitted = errors.New(`response already committed`)
)
