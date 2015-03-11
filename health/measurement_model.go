package health

import "runtime"

type measurements struct {
	NumCPU          int    `json:"num_cpu"`
	NumGoroutine    int    `json:"num_goroutine"`
	MemoryAllocated uint64 `json:"memory_allocated"`
}

func takeMeasurements() measurements {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	return measurements{
		NumCPU:          runtime.NumCPU(),
		NumGoroutine:    runtime.NumGoroutine(),
		MemoryAllocated: memStats.Alloc,
	}
}
