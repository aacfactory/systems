package commons

import (
	"os"
	"path/filepath"
)

func GetEnv(key string, default_ string, combineWith ...string) string {
	value := os.Getenv(key)
	if value == "" {
		value = default_
	}
	switch len(combineWith) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith)+1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
}

func HostProc(combineWith ...string) string {
	return GetEnv("HOST_PROC", "/proc", combineWith...)
}
