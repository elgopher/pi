// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package internal

import (
	"os"

	"github.com/shirou/gopsutil/v4/process"
)

var currentProcess *process.Process

func init() {
	var err error
	currentProcess, err = process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("Error creating process object: " + err.Error())
	}
}
