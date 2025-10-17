package oauth2

import "github.com/admpub/goth"

var constructors = map[string]func(*Account) goth.Provider{}

func Register(name string, constructor func(*Account) goth.Provider) {
	constructors[name] = constructor
}

func ConstructorNames() []string {
	var names []string
	for name := range constructors {
		names = append(names, name)
	}
	return names
}
