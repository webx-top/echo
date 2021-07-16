package embed

import "github.com/webx-top/echo"

type EmbedConfig struct {
	Index    string
	Prefix   string
	FilePath func(echo.Context) (string, error)
}

var DefaultEmbedConfig = EmbedConfig{
	Index: "index.html",
	FilePath: func(c echo.Context) (string, error) {
		return c.Param(`*`), nil
	},
}
