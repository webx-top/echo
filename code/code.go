package code

import (
	"net/http"
)

type (
	Code int
	TextHTTPCode struct {
		Text string
		HTTPCode int
	}
	CodeMap map[Code]TextHTTPCode
)

const (
	RequestTimeout Code = -9 //提交超时
	AbnormalResponse Code = -8 //响应异常
	OperationTimeout Code = -7 //操作超时
	Unsupported Code = -6 //不支持的操作
	RepeatOperation Code = -5 //重复操作
	DataNotFound Code = -4 //数据未找到
	UserNotFound Code = -3 //用户未找到
	NonPrivileged Code = -2 //无权限
	Unauthenticated Code = -1 //未登录
	Failure Code = 0 //操作失败
	Success Code = 1 //操作成功
)


// CodeDict 状态码字典
var CodeDict = CodeMap {
	RequestTimeout : {"RequestTimeout", http.StatusOK},
	AbnormalResponse : {"AbnormalResponse", http.StatusOK},
	OperationTimeout : {"OperationTimeout", http.StatusOK},
	Unsupported : {"Unsupported", http.StatusOK},
	RepeatOperation : {"RepeatOperation", http.StatusOK},
	DataNotFound : {"DataNotFound", http.StatusOK},
	UserNotFound : {"UserNotFound", http.StatusOK},
	NonPrivileged : {"NonPrivileged", http.StatusOK},
	Unauthenticated : {"Unauthenticated", http.StatusOK},
	Failure : {"Failure", http.StatusOK},
	Success : {"Success", http.StatusOK},
}

func (c Code) String() string {
	if v, y := CodeDict[c]; y {
		return v.Text
	}
	return `Undefined`
}

// Int 返回int类型的自定义状态码
func (c Code) Int() int {
	return int(c)
}

// HTTPCode 返回HTTP状态码
func (c Code) HTTPCode() int {
	if v, y := CodeDict[c]; y {
		return v.HTTPCode
	}
	return http.StatusOK
}

func (s CodeMap) Get(code int) TextHTTPCode {
	v, _ := s[Code(code)]
	return v
}
