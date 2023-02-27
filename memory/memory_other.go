//go:build !darwin && !linux && !windows

package memory

import "fmt"

func Stats() (mem *Memory, err error) {
	err = fmt.Errorf("unsupported os")
	return
}
