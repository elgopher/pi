// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package internal

import (
	"runtime"
)

var memStats = &runtime.MemStats{}

func GetMemoryMB() int {
	mem, err := currentProcess.MemoryInfo()
	if err != nil {
		panic(err)
	}
	return int(float64(mem.RSS) / 1024 / 1024)
}

func GetAllocatedObjectsCount() uint64 {
	prev := memStats.Mallocs
	runtime.ReadMemStats(memStats)
	return memStats.Mallocs - prev
}
