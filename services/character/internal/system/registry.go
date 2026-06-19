package system

import "errors"

var ErrUnknownSystem = errors.New("unknown system")

var registry = map[string]Character{}

func Register(name string, c Character) {
	registry[name] = c
}

func Get(name string) (Character, error) {
	c, ok := registry[name]
	if !ok {
		return nil, ErrUnknownSystem
	}
	return c, nil
}
