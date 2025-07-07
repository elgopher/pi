// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

const windowStateFileName = ".window"

type windowState struct {
	lastX, lastY, lastW, lastH int
	monitor                    int
}

func (s *windowState) store() {
	if RememberWindow {
		x, y := ebiten.WindowPosition()
		w, h := ebiten.WindowSize()

		var monitors []*ebiten.MonitorType
		monitors = ebiten.AppendMonitors(monitors[:0])
		currentMonitor := ebiten.Monitor()
		monitor := slices.Index(monitors, currentMonitor)

		if x != s.lastX || y != s.lastY || w != s.lastW || h != s.lastH || monitor != s.monitor {
			s.lastX, s.lastY = x, y
			s.lastW, s.lastH = w, h

			file := []byte(fmt.Sprintf("%d %d %d %d %d", x, y, w, h, monitor))

			if err := os.WriteFile(windowStateFileName, file, 0644); err != nil {
				log.Printf("Failed to save window state: %v", err)
				return
			}
		}
	}
}

func (s *windowState) restore() {
	if RememberWindow {
		state, err := os.ReadFile(windowStateFileName)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Failed to restore window state: %v", err)
			}
			return
		}

		buf := bytes.NewBuffer(state)
		_, err = fmt.Fscanf(buf, "%d %d %d %d %d", &s.lastX, &s.lastY, &s.lastW, &s.lastH, &s.monitor)
		if err != nil {
			log.Printf("Failed to restore window state: %v", err)
			return
		}

		ebiten.SetWindowPosition(s.lastX, s.lastY)
		ebiten.SetWindowSize(s.lastW, s.lastH)

		var monitors []*ebiten.MonitorType
		monitors = ebiten.AppendMonitors(monitors[:0])
		if s.monitor < len(monitors) {
			ebiten.SetMonitor(monitors[s.monitor])
		}
	}
}
