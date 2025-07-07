// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package internal

import (
	"math"
	"time"
)

func GetCPU() int {
	return int(math.Round(cpuUsageCollector.get() * 100))
}

var cpuUsageCollector = &cpuUsage{}

type cpuUsage struct {
	prevMeasureTime time.Time
	prevCpuTime     float64
}

func (u *cpuUsage) get() float64 {
	times, err := currentProcess.Times()
	if err != nil {
		panic(err)
	}

	currentCpuTime := times.User + times.System

	now := time.Now()
	realTimeElapsed := float64(now.Sub(u.prevMeasureTime)) / float64(time.Second)
	u.prevMeasureTime = now
	usage := (currentCpuTime - u.prevCpuTime) / realTimeElapsed
	u.prevCpuTime = currentCpuTime

	return usage
}
