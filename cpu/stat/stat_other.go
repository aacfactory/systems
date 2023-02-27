//go:build !linux && !windows

package stat

import "fmt"

func Times() (times []TimesStat, err error) {
	err = fmt.Errorf("not implements")
	return
}
