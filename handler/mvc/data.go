package mvc

type Data interface {
	Assign(key string, val interface{})
	Assignx(values *map[string]interface{})
	SetTmplFuncs()
	Set(code int, args ...interface{})
}
