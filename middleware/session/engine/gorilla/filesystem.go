package session

import (
	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
)

type FilesystemStore interface {
	Store
}

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(path string, keyPairs ...[]byte) FilesystemStore {
	return &filesystemStore{sessions.NewFilesystemStore(path, keyPairs...)}
}

type filesystemStore struct {
	*sessions.FilesystemStore
}

func (c *filesystemStore) Options(options echo.SessionOptions) {
	c.FilesystemStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
