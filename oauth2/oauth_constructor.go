package oauth2

import (
	"github.com/admpub/goth"
)

var constructors = map[string]func(*Config) goth.Provider{}

// ConstructorList returns the names of all registered constructor
func ConstructorList() []string {
	list := make([]string, 0, len(constructors))
	for k := range constructors {
		list = append(list, k)
	}
	return list
}

func Register(name string, newi func(cfg *Config) goth.Provider) {
	constructors[name] = newi
}

func GetConstructor(name string) func(cfg *Config) goth.Provider {
	return constructors[name]
}
