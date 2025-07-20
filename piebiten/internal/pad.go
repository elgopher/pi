// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/elgopher/pi/pipad"
)

var directionalButtons = map[ebiten.StandardGamepadButton]struct{}{
	ebiten.StandardGamepadButtonLeftLeft:   {},
	ebiten.StandardGamepadButtonLeftRight:  {},
	ebiten.StandardGamepadButtonLeftBottom: {},
	ebiten.StandardGamepadButtonLeftTop:    {},
}

var nonDirectionalButtonsMapping = map[ebiten.StandardGamepadButton]pipad.Button{
	ebiten.StandardGamepadButtonRightLeft:        pipad.X,
	ebiten.StandardGamepadButtonRightRight:       pipad.B,
	ebiten.StandardGamepadButtonRightBottom:      pipad.A,
	ebiten.StandardGamepadButtonRightTop:         pipad.Y,
	ebiten.StandardGamepadButtonFrontTopRight:    "RB",
	ebiten.StandardGamepadButtonFrontBottomRight: "RT", // right trigger
	ebiten.StandardGamepadButtonFrontTopLeft:     "LB",
	ebiten.StandardGamepadButtonFrontBottomLeft:  "LT", // left trigger
	ebiten.StandardGamepadButtonCenterLeft:       "Back",
	ebiten.StandardGamepadButtonCenterRight:      "Start",
	ebiten.StandardGamepadButtonLeftStick:        "LSB", // left stick button
	ebiten.StandardGamepadButtonRightStick:       "RSB", // right stick button
}

type ebitenGamepads struct {
	pads                   [16]pad
	gamepadIDs             []ebiten.GamepadID
	standardGamepadButtons []ebiten.StandardGamepadButton
}

func (g *ebitenGamepads) update() {
	g.publishDisconnectedGamepadEvents()
	g.gamepadIDs = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIDs[:0])
	g.publishConnectedGamepadEvents()
	g.gamepadIDs = ebiten.AppendGamepadIDs(g.gamepadIDs[:0])
	for _, gamepadID := range g.gamepadIDs {
		player := int(gamepadID)
		g.publishDownEventsForNonDirectionalGamepadButtons(gamepadID, player)
		g.publishUpEventsForNonDirectionalGamepadButtons(gamepadID, player)
		g.publishEventsForDirectionalGamepadButtons(gamepadID, player)
	}
}

func (g *ebitenGamepads) publishDisconnectedGamepadEvents() {
	for _, id := range g.gamepadIDs {
		if inpututil.IsGamepadJustDisconnected(id) {
			pipad.ConnectionTarget().Publish(pipad.EventConnection{
				Type:   pipad.EventDisconnect,
				Player: int(id),
			})
		}
	}
}

func (g *ebitenGamepads) publishConnectedGamepadEvents() {
	for _, id := range g.gamepadIDs {
		pipad.ConnectionTarget().Publish(pipad.EventConnection{
			Type:   pipad.EventConnect,
			Player: int(id),
		})
	}
}

func (g *ebitenGamepads) publishDownEventsForNonDirectionalGamepadButtons(gamepadID ebiten.GamepadID, player int) {
	g.standardGamepadButtons = inpututil.AppendJustPressedStandardGamepadButtons(gamepadID, g.standardGamepadButtons[:0])
	for _, button := range g.standardGamepadButtons {
		if _, isDirectional := directionalButtons[button]; !isDirectional {
			if mapping := nonDirectionalButtonsMapping[button]; mapping != "" {
				pipad.ButtonTarget().Publish(
					pipad.EventButton{Type: pipad.EventDown, Button: mapping, Player: player},
				)
			}
		}
	}
}

func (g *ebitenGamepads) publishUpEventsForNonDirectionalGamepadButtons(gamepadID ebiten.GamepadID, player int) {
	g.standardGamepadButtons = inpututil.AppendJustReleasedStandardGamepadButtons(gamepadID, g.standardGamepadButtons[:0])
	for _, button := range g.standardGamepadButtons {
		if _, isDirectional := directionalButtons[button]; !isDirectional {
			if mapping := nonDirectionalButtonsMapping[button]; mapping != "" {
				pipad.ButtonTarget().Publish(
					pipad.EventButton{Type: pipad.EventUp, Button: mapping, Player: player},
				)
			}
		}
	}
}

func (g *ebitenGamepads) publishEventsForDirectionalGamepadButtons(gamepadID ebiten.GamepadID, player int) {
	prevPad := g.pads[player]
	newPad := g.pads[player]
	newPad.handleInput(gamepadID)
	newPad.publishEvents(prevPad, player)
	g.pads[player] = newPad
}

type pad struct {
	left, right, top, bottom padButton
}

func (p *pad) handleInput(gamepadID ebiten.GamepadID) {
	p.handleAxisInput(gamepadID)
	p.handleButtonsInput(gamepadID)
}

func (p *pad) handleAxisInput(gamepadID ebiten.GamepadID) {
	horizontal := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal)
	vertical := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickVertical)
	p.left.axisDown = btoi(horizontal < -0.5)
	p.right.axisDown = btoi(horizontal > 0.5)
	p.top.axisDown = btoi(vertical < -0.5)
	p.bottom.axisDown = btoi(vertical > 0.5)
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (p *pad) handleButtonsInput(gamepadID ebiten.GamepadID) {
	p.left.dpadDown = btoi(ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftLeft))
	p.right.dpadDown = btoi(ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftRight))
	p.top.dpadDown = btoi(ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftTop))
	p.bottom.dpadDown = btoi(ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftBottom))
}

func (p *pad) publishEvents(prev pad, player int) {
	p.left.publishEvent(prev.left, pipad.Left, player)
	p.right.publishEvent(prev.right, pipad.Right, player)
	p.top.publishEvent(prev.top, pipad.Top, player)
	p.bottom.publishEvent(prev.bottom, pipad.Bottom, player)
}

type padButton struct {
	axisDown int // 0 or 1 (1 is true)
	dpadDown int // 0 or 1 (1 is true)
}

func (b padButton) publishEvent(prev padButton, button pipad.Button, player int) {
	prevDown := prev.axisDown + prev.dpadDown
	down := b.axisDown + b.dpadDown
	if prevDown > 0 && down == 0 {
		pipad.ButtonTarget().Publish(
			pipad.EventButton{Type: pipad.EventUp, Button: button, Player: player},
		)
		return
	}
	if prevDown == 0 && down > 0 {
		pipad.ButtonTarget().Publish(
			pipad.EventButton{Type: pipad.EventDown, Button: button, Player: player},
		)
		return
	}
}
