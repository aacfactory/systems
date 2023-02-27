//go:build linux

package memory

import (
	"fmt"
	"github.com/aacfactory/systems/internal/commons"
	"math"
	"os"
	"strconv"
	"strings"
)

func Stats() (mem *Memory, err error) {
	filename := commons.HostProc("meminfo")
	lines, readErr := commons.ReadLines(filename)
	if readErr != nil {
		err = readErr
		return
	}
	hasAvailable := false
	buffers := uint64(0)
	cached := uint64(0)
	sReclaimable := uint64(0)
	hasSRReclaimable := false
	activeFile := uint64(0)
	hasActiveFile := false
	inactiveFile := uint64(0)
	hasInactiveFile := false
	mem = &Memory{}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)
		switch key {
		case "MemTotal":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse MemTotal failed")
				return
			}
			mem.Total = n * 1024
			break
		case "MemFree":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse MemFree failed")
				return
			}
			mem.Free = n * 1024
			break
		case "MemAvailable":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse MemAvailable failed")
				return
			}
			hasAvailable = true
			mem.Available = n * 1024
			break
		case "Buffers":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse Buffers failed")
				return
			}
			buffers = n * 1024
			break
		case "Cached":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse Cached failed")
				return
			}
			cached = n * 1024
			break
		case "SReclaimable":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse SReclaimable failed")
				return
			}
			sReclaimable = n * 1024
			hasSRReclaimable = true
			break
		case "Active(file)":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse Active(file) failed")
				return
			}
			activeFile = n * 1024
			hasActiveFile = true
			break
		case "Inactive(file)":
			n, parseErr := strconv.ParseUint(value, 10, 64)
			if parseErr != nil {
				err = fmt.Errorf("parse Inactive(file) failed")
				return
			}
			inactiveFile = n * 1024
			hasInactiveFile = true
			break
		default:
			break
		}
	}
	cached += sReclaimable
	if !hasAvailable {
		if hasActiveFile && hasInactiveFile && hasSRReclaimable {
			fn := commons.HostProc("zoneinfo")
			lines, readErr = commons.ReadLines(fn)
			if readErr != nil {
				err = readErr
				mem.Available = mem.Free + cached
			} else {
				pageSize := uint64(os.Getpagesize())
				watermarkLow := uint64(0)
				for _, line := range lines {
					fields := strings.Fields(line)

					if strings.HasPrefix(fields[0], "low") {
						lowValue, err := strconv.ParseUint(fields[1], 10, 64)
						if err != nil {
							lowValue = 0
						}
						watermarkLow += lowValue
					}
				}
				watermarkLow *= pageSize
				availMemory := mem.Free - watermarkLow
				pageCache := activeFile + inactiveFile
				pageCache -= uint64(math.Min(float64(pageCache/2), float64(watermarkLow)))
				availMemory += pageCache
				availMemory += sReclaimable - uint64(math.Min(float64(sReclaimable/2.0), float64(watermarkLow)))
				if availMemory < 0 {
					availMemory = 0
				}
			}
		} else {
			mem.Available = cached + mem.Free
		}
	}
	mem.Used = mem.Total - mem.Free - buffers - cached
	mem.UsedPercent = float64(mem.Used) / float64(mem.Total) * 100.0
	fillRuntime(mem)
	return
}
