package oauth2

import "github.com/admpub/goth"

var constructors = map[string]func(*Account) goth.Provider{}

// Register adds a new OAuth2 provider constructor to the registry with the given name.
func Register(name string, constructor func(*Account) goth.Provider) {
	constructors[name] = constructor
}

// GetConstructor returns the provider constructor function for the given name.
func GetConstructor(name string) func(*Account) goth.Provider {
	return constructors[name]
}

// ConstructorNames returns a slice containing all registered constructor names.
func ConstructorNames() []string {
	names := make([]string, 0, len(constructors))
	for name := range constructors {
		names = append(names, name)
	}
	return names
}
