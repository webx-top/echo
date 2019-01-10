package host

// Configer 配置接口
type Configer interface {
	Init(conf interface{}) Configer
	Set(data string) error
	Get(recvFn func(interface{}))
	Config() interface{}
	IsValid() bool
	String() string
	Template() string                     //获取表单模板名称
	SetTemplate(tmplFile string) Configer //获取设置表单模板
}
