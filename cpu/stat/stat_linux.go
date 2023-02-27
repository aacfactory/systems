//go:build linux

package stat

import (
	"errors"
	"github.com/aacfactory/systems/internal/commons"
	"strconv"
	"strings"
)

func Times() (times []TimesStat, err error) {
	filename := commons.HostProc("stat")
	lines := []string{}
	statlines, readErr := commons.ReadLines(filename)
	if readErr != nil || len(statlines) < 2 {
		return []TimesStat{}, nil
	}
	for _, line := range statlines[1:] {
		if !strings.HasPrefix(line, "cpu") {
			break
		}
		lines = append(lines, line)
	}
	times = make([]TimesStat, 0, len(lines))
	for _, line := range lines {
		ct, parseErr := parseStatLine(line)
		if parseErr != nil {
			continue
		}
		times = append(times, *ct)
	}
	return
}

func parseStatLine(line string) (*TimesStat, error) {
	fields := strings.Fields(line)

	if len(fields) == 0 {
		return nil, errors.New("stat does not contain cpu info")
	}
	cpu := fields[0]
	if !strings.HasPrefix(cpu, "cpu") {
		return nil, errors.New("not contain cpu")
	}
	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	nice, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}
	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	irq, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, err
	}
	softirq, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, err
	}

	ct := &TimesStat{
		CPU:     cpu,
		User:    user / float64(100),
		Nice:    nice / float64(100),
		System:  system / float64(100),
		Idle:    idle / float64(100),
		Iowait:  iowait / float64(100),
		Irq:     irq / float64(100),
		Softirq: softirq / float64(100),
	}
	if len(fields) > 8 {
		steal, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, err
		}
		ct.Steal = steal / float64(100)
	}
	if len(fields) > 9 {
		guest, err := strconv.ParseFloat(fields[9], 64)
		if err != nil {
			return nil, err
		}
		ct.Guest = guest / float64(100)
	}
	if len(fields) > 10 {
		guestNice, err := strconv.ParseFloat(fields[10], 64)
		if err != nil {
			return nil, err
		}
		ct.GuestNice = guestNice / float64(100)
	}

	return ct, nil
}
