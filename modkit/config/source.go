package config

import "os"

// Source resolves a config key to a raw string value.
type Source interface {
	Lookup(key string) (value string, ok bool)
}

type envSource struct{}

func (envSource) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}
