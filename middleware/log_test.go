package middleware

import (
	"testing"
)

func TestLogger(t *testing.T) {
	for _, c := range terminalColors {
		c.PrintlnFunc()(`TEST`)
	}
}
