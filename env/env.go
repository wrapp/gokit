package env

import (
	"os"
	"strconv"
	"strings"
)

// Default returns the environment variable against the provided `key`. If there is no such
// variable then the default value be returned provided in `def`.
func Default(key, def string) string {
	e := os.Getenv(key)
	if e == "" {
		return def
	}
	return e
}

// Get returns the environment variable agains the provided `key`. If there is no such
// variable then empty string is returned.
func Get(key string) string {
	return Default(key, "")
}

// DefaultInt returns the environment variable against the provided `key` as an integer. If
// there is no such variable or the value cannot be converted into an integer then a default
// is returned provided in `def`.
func DefaultInt(key string, def int) int {
	env := Get(key)
	if env == "" {
		return def
	}
	i, err := strconv.ParseInt(env, 10, 32)
	if err != nil {
		return def
	}
	return int(i)
}

// Bool returns the environment variable against the provided `key`. If there is no such
// variable then `false`is returned. The values which are considered `true` are '1', 'true',
// 'yes' and 'on'. These values are case-insensetive.
func Bool(key string) bool {
	env := strings.ToLower(os.Getenv(key))
	if env == "1" || env == "true" || env == "yes" || env == "on" {
		return true
	}
	return false
}

// ServiceName returns the value of environment variable `SERVICE_NAME`.
func ServiceName() string {
	return Get("SERVICE_NAME")
}
