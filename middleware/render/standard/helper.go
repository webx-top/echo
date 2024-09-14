package standard

import (
	htmlTpl "html/template"
	"regexp"

	"github.com/webx-top/echo"
)

var nlspInBlock = regexp.MustCompile("^[ \t\r]*\n|\n[ \t\r]*$")

func trimSpaceInBlock(innerData string) string {
	return nlspInBlock.ReplaceAllString(innerData, ``)
}

func setFunc(tplInf *tplInfo, funcMap htmlTpl.FuncMap) htmlTpl.FuncMap {
	if funcMap == nil {
		funcMap = htmlTpl.FuncMap{}
	}
	funcMap["hasBlock"] = func(blocks ...string) bool {
		for _, blockName := range blocks {
			if _, ok := tplInf.Blocks[blockName]; !ok {
				return false
			}
		}
		return true
	}
	funcMap["hasAnyBlock"] = func(blocks ...string) bool {
		for _, blockName := range blocks {
			if _, ok := tplInf.Blocks[blockName]; ok {
				return true
			}
		}
		return false
	}
	return funcMap
}

func parseError(err error, sourceContent string) *echo.PanicError {
	return echo.ParseTemplateError(err, sourceContent)
}
