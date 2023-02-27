//go:build darwin

package memory

import (
	"context"
	"fmt"
	"github.com/aacfactory/systems/internal/commons"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

func Stats() (mem *Memory, err error) {
	total, readMemSizeErr := unix.SysctlUint64("hw.memsize")
	if readMemSizeErr != nil {
		err = readMemSizeErr
		return
	}
	vmBytes, callErr := commons.ExecCommand(context.TODO(), "vm_stat")
	if callErr != nil {
		err = callErr
		return
	}

	inactive := uint64(0)

	mem = &Memory{
		Total:       total,
		Free:        0,
		Available:   0,
		Used:        0,
		UsedPercent: 0,
		Heap:        nil,
		Stack:       nil,
		GC:          nil,
	}
	lines := strings.Split(string(vmBytes), "\n")
	pagesize := uint64(unix.Getpagesize())
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.Trim(fields[1], " .")
		switch key {
		case "Pages free":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse Pages free failed")
				return
			}
			mem.Free = n * pagesize
			break
		case "Pages inactive":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse Inactive(file) failed")
				return
			}
			inactive = n * pagesize
			break
		default:
			break
		}
	}

	mem.Available = mem.Free + inactive
	mem.Used = mem.Total - mem.Available
	mem.UsedPercent = 100 * float64(mem.Used) / float64(mem.Total)
	fillRuntime(mem)
	return
}
