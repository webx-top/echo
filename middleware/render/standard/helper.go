package standard

import (
	"regexp"

	"github.com/webx-top/echo"
)

var nlspInBlock = regexp.MustCompile("^[ \t\r]*\n|\n[ \t\r]*$")

func trimSpaceInBlock(innerData string) string {
	return nlspInBlock.ReplaceAllString(innerData, ``)
}

func parseError(err error, sourceContent string) *echo.PanicError {
	return echo.ParseTemplateError(err, sourceContent)
}
