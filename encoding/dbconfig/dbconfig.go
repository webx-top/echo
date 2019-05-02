package dbconfig

var DefaultCharset = `utf8`

type Config struct {
	Engine  string
	User    string
	Pass    string
	Name    string
	Host    string
	Port    string
	Charset string
	Prefix  string
	Options map[string]string
}

func (c *Config) String() string {
	if c.Options == nil {
		c.Options = map[string]string{}
	}
	encoder := EncoderGet(c.Engine)
	if encoder == nil {
		panic(c.Engine + ` is not supported.`)
	}
	return encoder(c)
}
