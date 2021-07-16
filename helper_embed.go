// +build go1.16

package echo

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

func NewFileSystems() FileSystems {
	return FileSystems{}
}

type FileSystems []fs.FS

func (f FileSystems) Open(name string) (file fs.File, err error) {
	for _, fileSystem := range f {
		file, err = fileSystem.Open(name)
		if err == nil || !errors.Is(err, fs.ErrNotExist) {
			return
		}
	}
	return
}

func (f FileSystems) Size() int {
	return len(f)
}

func (f FileSystems) IsEmpty() bool {
	return f.Size() == 0
}

func (f *FileSystems) Register(fileSystem fs.FS) {
	*f = append(*f, fileSystem)
}

type EmbedConfig struct {
	Index    string
	Prefix   string
	FilePath func(Context) (string, error)
}

var DefaultEmbedConfig = EmbedConfig{
	Index: "index.html",
	FilePath: func(c Context) (string, error) {
		return c.Param(`*`), nil
	},
}

// EmbedFile
// e.Get(`/*`, EmbedFile(customFS))
func EmbedFile(fs FileSystems, configs ...EmbedConfig) func(c Context) error {
	config := DefaultEmbedConfig
	if len(configs) > 0 {
		config = configs[0]
		if len(config.Index) == 0 {
			config.Index = DefaultEmbedConfig.Index
		}
		if len(config.Prefix) > 0 {
			config.Prefix = strings.TrimPrefix(config.Prefix, `/`)
		}
		if len(config.Prefix) > 0 {
			if !strings.HasSuffix(config.Prefix, `/`) {
				config.Prefix += `/`
			}
		}
		if config.FilePath == nil {
			config.FilePath = DefaultEmbedConfig.FilePath
		}
	}
	return func(c Context) error {
		file, err := config.FilePath(c)
		if err != nil {
			return err
		}
		if len(file) == 0 {
			file = config.Index
		}
		if len(config.Prefix) > 0 {
			file = config.Prefix + file
		}
		f, err := fs.Open(file)
		if err != nil {
			return ErrNotFound
		}
		defer func() {
			if f != nil {
				f.Close()
			}
		}()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			f.Close()

			file = filepath.Join(file, DefaultEmbedConfig.Index)
			if f, err = fs.Open(file); err != nil {
				return ErrNotFound
			}

			if fi, err = f.Stat(); err != nil {
				return err
			}
		}
		return c.ServeContent(f, fi.Name(), fi.ModTime())
	}
}
