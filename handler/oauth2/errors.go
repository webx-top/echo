package oauth2

import "errors"

var (
	ErrStateTokenMismatch = errors.New("state token mismatch")
	ErrSessionDismatched  = errors.New("could not find a matching session for this request")
	ErrMustSelectProvider = errors.New("you must select a provider")

	// Unpack Value
	ErrIPAddressDismatched = errors.New(`IP address does not match`)
	ErrUserAgentDismatched = errors.New(`UserAgent does not match`)
	ErrDataExpired         = errors.New(`data has expired`)
)
