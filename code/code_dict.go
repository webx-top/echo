package code

import "net/http"

// CodeDict 状态码字典
var CodeDict = CodeMap{

	// - 系统状态

	SystemUnauthorized: {"SystemUnauthorized", http.StatusServiceUnavailable},
	SystemNotInstalled: {"SystemNotInstalled", http.StatusServiceUnavailable},

	// - 操作状态

	OperationProcessing: {"OperationProcessing", http.StatusProcessing},
	FrequencyTooFast:    {"FrequencyTooFast", http.StatusTooManyRequests},
	RequestFailure:      {"RequestFailure", http.StatusBadRequest},
	RequestTimeout:      {"RequestTimeout", http.StatusRequestTimeout},
	AbnormalResponse:    {"AbnormalResponse", http.StatusInternalServerError},
	OperationTimeout:    {"OperationTimeout", http.StatusRequestTimeout},
	Unsupported:         {"Unsupported", http.StatusNotImplemented},
	RepeatOperation:     {"RepeatOperation", http.StatusBadRequest},

	// - 数据状态

	InvalidAppID: {"InvalidAppID", http.StatusBadRequest},
	InvalidToken: {"InvalidToken", http.StatusBadRequest},

	DataNotChanged:      {"DataNotChanged", http.StatusNotModified},
	DataSizeTooBig:      {"DataSizeTooBig", http.StatusRequestEntityTooLarge},
	DataAlreadyExists:   {"DataAlreadyExists", http.StatusBadRequest},
	DataFormatIncorrect: {"DataFormatIncorrect", http.StatusBadRequest},
	DataStatusIncorrect: {"DataStatusIncorrect", http.StatusBadRequest},
	DataHasExpired:      {"DataHasExpired", http.StatusBadRequest},
	DataProcessing:      {"DataProcessing", http.StatusProcessing},
	DataUnavailable:     {"DataUnavailable", http.StatusBadRequest},
	InvalidType:         {"InvalidType", http.StatusBadRequest},
	InvalidSignature:    {"InvalidSignature", http.StatusBadRequest},
	InvalidParameter:    {"InvalidParameter", http.StatusBadRequest},
	DataNotFound:        {"DataNotFound", http.StatusNotFound},

	// - 验证码

	CaptchaCodeRequired: {"CaptchaCodeRequired", http.StatusBadRequest},
	CaptchaIdMissing:    {"CaptchaIdMissing", http.StatusBadRequest},
	CaptchaError:        {"CaptchaError", http.StatusBadRequest},

	// - 用户状态

	BalanceNoEnough: {"BalanceNoEnough", http.StatusPreconditionFailed},
	UserDisabled:    {"UserDisabled", http.StatusBadRequest},
	UserNotFound:    {"UserNotFound", http.StatusBadRequest},
	NonPrivileged:   {"NonPrivileged", http.StatusForbidden},
	Unauthenticated: {"Unauthenticated", http.StatusUnauthorized},

	// - 通用

	Failure: {"Failure", http.StatusInternalServerError},
	Success: {"Success", http.StatusOK},
}
