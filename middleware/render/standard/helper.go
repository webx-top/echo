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

func replaceByMatchedIndex(content string, matches [][]int, replaced *string) func(k int, v []int, newInnerStr ...string) {
	endK := len(matches) - 1
	var lastEndIdx int
	return func(k int, v []int, newInnerStr ...string) {
		if len(newInnerStr) > 0 {
			if k == 0 {
				*replaced = content[0:v[0]] + newInnerStr[0]
				if k == endK {
					*replaced += content[v[1]:]
				}
			} else if k == endK {
				*replaced += newInnerStr[0] + content[v[1]:]
			} else {
				*replaced += content[lastEndIdx:v[0]] + newInnerStr[0]
			}
		} else {
			if k == 0 {
				if k == endK {
					*replaced = content
				} else {
					*replaced = content[0:v[1]]
				}
			} else if k == endK {
				*replaced += content[v[0]:]
			} else {
				*replaced += content[lastEndIdx:v[1]]
			}
		}
		lastEndIdx = v[1]
	}
}

func getMatchedByIndex(content string, v []int, recv ...*string) {
	recvNum := len(recv)
	matchIdx := 0
	matchNum := len(v)
	for idx := 0; idx < recvNum; idx++ {
		if recv[idx] != nil && v[matchIdx] > -1 {
			endIdx := matchIdx + 1
			if endIdx >= matchNum {
				return
			}
			*(recv[idx]) = content[v[matchIdx]:v[endIdx]]
		}
		matchIdx += 2
	}
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
