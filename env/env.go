package env

import "os"

func GetDefault(key, def string) string {
	e := os.Getenv(key)
	if e == "" {
		return def
	}
	return e
}

func Get(key string) string {
	return GetDefault(key, "")
}
