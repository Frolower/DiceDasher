package system

import "errors"

var ErrUnknownSystem = errors.New("unknown system")

var registry = map[string]Resolver{}

func Register(name string, r Resolver) {
	registry[name] = r
}

func Get(name string) (Resolver, error) {
	r, ok := registry[name]
	if !ok {
		return nil, ErrUnknownSystem
	}
	return r, nil
}
