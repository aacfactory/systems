package cpu

import (
	"fmt"
	"github.com/aacfactory/systems/cpu/stat"
	"math"
	"runtime"
	"time"
)

type Core struct {
	No        string  `json:"no"`
	Occupancy float64 `json:"occupancy"`
}

type CPU []Core

func (c CPU) Max() (core Core) {
	core = c[0]
	for _, e := range c {
		if core.Occupancy > e.Occupancy {
			core = e
		}
	}
	return
}

func (c CPU) Min() (core Core) {
	core = c[0]
	for _, e := range c {
		if core.Occupancy < e.Occupancy {
			core = e
		}
	}
	return
}

func (c CPU) AVG() (occupancy float64) {
	occupancy = float64(0)
	for _, core := range c {
		occupancy += core.Occupancy
	}
	occupancy = occupancy / float64(len(c))
	return
}

func Occupancy() (cpu CPU, err error) {
	cpuTimes1, time1Err := stat.Times()
	if time1Err != nil {
		err = time1Err
		return
	}
	time.Sleep(10 * time.Millisecond)
	cpuTimes2, time2Err := stat.Times()
	if time2Err != nil {
		err = time2Err
		return
	}
	if len(cpuTimes1) != len(cpuTimes2) {
		err = fmt.Errorf("received two CPU counts: %d != %d", len(cpuTimes1), len(cpuTimes2))
		return
	}
	cores := make([]Core, len(cpuTimes1))
	for i, t1 := range cpuTimes1 {
		var to stat.TimesStat
		matched := false
		for _, t2 := range cpuTimes2 {
			if t1.CPU == t2.CPU {
				to = t2
				matched = true
				break
			}
		}
		if !matched {
			err = fmt.Errorf("received two CPU has not same no")
			return
		}
		cores[i] = Core{
			No:        t1.CPU,
			Occupancy: calculateBusy(t1, to),
		}
	}
	cpu = cores
	return
}

func getAllBusy(t stat.TimesStat) (float64, float64) {
	tot := t.Total()
	if runtime.GOOS == "linux" {
		tot -= t.Guest
		tot -= t.GuestNice
	}
	busy := tot - t.Idle - t.Iowait
	return tot, busy
}

func calculateBusy(t1, t2 stat.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)
	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}
