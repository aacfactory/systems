package memory

import "runtime"

type Memory struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Heap        *Heap   `json:"heap"`
	Stack       *Stack  `json:"stack"`
	GC          *GC     `json:"gc"`
}

type Heap struct {
	Sys     uint64 `json:"sys"`
	Alloc   uint64 `json:"alloc"`
	Idle    uint64 `json:"idle"`
	Objects uint64 `json:"objects"`
}

type Stack struct {
	Sys uint64 `json:"sys"`
}

type GC struct {
	Enabled               bool   `json:"enabled"`
	Num                   uint32 `json:"num"`
	NumForced             uint32 `json:"numForced"`
	Sys                   uint64 `json:"sys"`
	Next                  uint64 `json:"next"`
	Last                  uint64 `json:"last"`
	PauseTotalNanoseconds uint64 `json:"pauseTotalNanoseconds"`
}

func fillRuntime(mem *Memory) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	heap := &Heap{
		Sys:     stats.HeapSys,
		Alloc:   stats.HeapAlloc,
		Idle:    stats.HeapIdle,
		Objects: stats.HeapObjects,
	}
	mem.Heap = heap
	stack := &Stack{
		Sys: stats.StackSys,
	}
	mem.Stack = stack
	gc := &GC{
		Enabled:               stats.EnableGC,
		Num:                   stats.NumGC,
		NumForced:             stats.NumForcedGC,
		Sys:                   stats.GCSys,
		Next:                  stats.NextGC,
		Last:                  stats.LastGC,
		PauseTotalNanoseconds: stats.PauseTotalNs,
	}
	mem.GC = gc
	return
}
