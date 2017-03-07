package env

import (
	"os"
	"strconv"
	"strings"
)

func Default(key, def string) string {
	e := os.Getenv(key)
	if e == "" {
		return def
	}
	return e
}

func Get(key string) string {
	return Default(key, "")
}

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

func Bool(key string) bool {
	env := strings.ToLower(os.Getenv(key))
	if env == "1" || env == "true" || env == "yes" || env == "on" {
		return true
	}
	return false
}

func ServiceName() string {
	return Get("SERVICE_NAME")
}
