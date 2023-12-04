package oauth2

import (
	"github.com/admpub/goth"
	"github.com/webx-top/echo"
)

type Account struct {
	On          bool // on / off
	Name        string
	Key         string
	Secret      string `json:"-" xml:"-"`
	Extra       echo.H
	LoginURL    string
	CallbackURL string
	Scopes      []string
	Constructor func(*Account) goth.Provider `json:"-" xml:"-"`
}

func (a *Account) SetConstructor(constructor func(*Account) goth.Provider) *Account {
	a.Constructor = constructor
	return a
}

func (a *Account) Instance() goth.Provider {
	return a.Constructor(a)
}
