package host

import "errors"

var (
	ErrHasBeenDisabled     = errors.New(`This module has been disabled`)
	ErrHasNotBeenInstalled = errors.New(`This module has not been installed`)
	ErrHasExpired          = errors.New(`This module has expired`)
)
