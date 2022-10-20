// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/vm"
)

var gameStoppedErr = errors.New("game stopped")

const tps = 30

func run() error {
	ebiten.SetTPS(tps)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(vm.ScreenWidth*scale(), vm.ScreenHeight*scale())
	ebiten.SetWindowSizeLimits(vm.ScreenWidth, vm.ScreenHeight, -1, -1)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowTitle("Pi Game")

	if err := ebiten.RunGame(&ebitengineGame{}); err != nil {
		if err == gameStoppedErr {
			return nil
		}

		return fmt.Errorf("running game using Ebitengine failed: %w", err)
	}

	return nil
}

func scale() int {
	return int(math.Round(ebiten.DeviceScaleFactor() * 3))
}

type ebitengineGame struct {
	screenDataRGBA     []byte // reused RGBA pixels
	screenChanged      bool
	shouldSkipNextDraw bool
}

func (e *ebitengineGame) Update() error {
	updateStartedTime := time.Now()

	updateTime()
	updateController()
	updateMouse()
	updateKeyDuration()

	if vm.Update != nil {
		vm.Update()
	}

	if vm.GameLoopStopped {
		return gameStoppedErr
	}

	// Ebitengine treats Draw differently than π. In π Draw must be executed at most 30 times per second.
	// That's why π runs Draw() from inside Ebitengine's Update().
	if vm.Draw != nil {
		if e.shouldSkipNextDraw {
			e.shouldSkipNextDraw = false
			return nil
		}

		vm.Draw()

		elapsed := time.Since(updateStartedTime)
		if elapsed.Seconds() > 1/float64(tps) {
			e.shouldSkipNextDraw = true
		}
	}

	e.screenChanged = true

	return nil
}

func (e *ebitengineGame) Draw(screen *ebiten.Image) {
	// Ebitengine executes Draw based on display frequency.
	// But the screen is changed at most 30 times per second.
	// That's why there is no need to write pixels more often
	// than 30 times per second.
	if e.screenChanged {
		e.writeScreenPixels(screen)
		e.screenChanged = false
	}
}

func (e *ebitengineGame) writeScreenPixels(screen *ebiten.Image) {
	if e.screenDataRGBA == nil || len(e.screenDataRGBA)/4 != len(vm.ScreenData) {
		e.screenDataRGBA = make([]byte, len(vm.ScreenData)*4)
	}

	offset := 0
	for _, col := range vm.ScreenData {
		rgb := vm.Palette[vm.DisplayPalette[col]]
		e.screenDataRGBA[offset] = rgb.R
		e.screenDataRGBA[offset+1] = rgb.G
		e.screenDataRGBA[offset+2] = rgb.B
		e.screenDataRGBA[offset+3] = 0xff
		offset += 4
	}

	screen.WritePixels(e.screenDataRGBA)
}

func (e *ebitengineGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return vm.ScreenWidth, vm.ScreenHeight
}

var lastTime time.Time

func updateTime() {
	now := time.Now()
	timePassed := now.Sub(lastTime)
	lastTime = now
	vm.TimeSeconds += float64(timePassed) / float64(time.Second)
}

func updateController() {
	for player := 0; player < 8; player++ {
		getController(player).update(player)
	}
}

func getController(player int) *controller {
	c := controller{&vm.Controllers[player]}
	return &c
}

type controller struct {
	*vm.Controller
}

func (c *controller) update(player int) {
	c.updateDirections(player)
	c.updateFireButtons(player)
}

func (c *controller) updateDirections(player int) {
	gamepadID := ebiten.GamepadID(player)

	axisX := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickHorizontal)
	axisY := ebiten.StandardGamepadAxisValue(gamepadID, ebiten.StandardGamepadAxisLeftStickVertical)

	if axisX < -0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftLeft) ||
		isKeyboardPressed(player, vm.ControllerLeft) {
		c.BtnDuration[vm.ControllerLeft] += 1
		c.BtnDuration[vm.ControllerRight] = 0
	} else if axisX > 0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftRight) ||
		isKeyboardPressed(player, vm.ControllerRight) {
		c.BtnDuration[vm.ControllerRight] += 1
		c.BtnDuration[vm.ControllerLeft] = 0
	} else {
		c.BtnDuration[vm.ControllerRight] = 0
		c.BtnDuration[vm.ControllerLeft] = 0
	}

	if axisY < -0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftTop) ||
		isKeyboardPressed(player, vm.ControllerUp) {
		c.BtnDuration[vm.ControllerUp] += 1
		c.BtnDuration[vm.ControllerDown] = 0
	} else if axisY > 0.5 ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonLeftBottom) ||
		isKeyboardPressed(player, vm.ControllerDown) {
		c.BtnDuration[vm.ControllerDown] += 1
		c.BtnDuration[vm.ControllerUp] = 0
	} else {
		c.BtnDuration[vm.ControllerUp] = 0
		c.BtnDuration[vm.ControllerDown] = 0
	}
}

func (c *controller) updateFireButtons(player int) {
	gamepadID := ebiten.GamepadID(player)

	if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightBottom) ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightTop) ||
		isKeyboardPressed(player, vm.ControllerO) {
		c.BtnDuration[vm.ControllerO] += 1
	} else {
		c.BtnDuration[vm.ControllerO] = 0
	}

	if ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightRight) ||
		ebiten.IsStandardGamepadButtonPressed(gamepadID, ebiten.StandardGamepadButtonRightLeft) ||
		isKeyboardPressed(player, vm.ControllerX) {
		c.BtnDuration[vm.ControllerX] += 1
	} else {
		c.BtnDuration[vm.ControllerX] = 0
	}
}

// first array is player, then π key, then slice of Ebitengine keys.
var keyboardMapping = [...][6][]ebiten.Key{
	// player0:
	{
		{ebiten.KeyLeft},                        // left
		{ebiten.KeyRight},                       // right
		{ebiten.KeyUp},                          // up
		{ebiten.KeyDown},                        // down
		{ebiten.KeyZ, ebiten.KeyC, ebiten.KeyN}, // o
		{ebiten.KeyX, ebiten.KeyV, ebiten.KeyM}, // x
	},
	// player1:
	{
		{ebiten.KeyS},         // left
		{ebiten.KeyF},         // right
		{ebiten.KeyE},         // up
		{ebiten.KeyD},         // down
		{ebiten.KeyShiftLeft}, // o
		{ebiten.KeyTab, ebiten.KeyW, ebiten.KeyQ, ebiten.KeyA}, // x
	},
}

func isKeyboardPressed(player int, button int) bool {
	if player >= len(keyboardMapping) {
		return false
	}

	keys := keyboardMapping[player][button]
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}

	return false
}

var mouseMapping = []ebiten.MouseButton{
	ebiten.MouseButtonLeft,
	ebiten.MouseButtonMiddle,
	ebiten.MouseButtonRight,
}

func updateMouse() {
	for i := 0; i < len(mouseMapping); i++ {
		button := mouseMapping[i]
		if ebiten.IsMouseButtonPressed(button) {
			vm.MouseBtnDuration[i] += 1
		} else {
			vm.MouseBtnDuration[i] = 0
		}
	}

	x, y := ebiten.CursorPosition()
	vm.MousePos.X = x
	vm.MousePos.Y = y
}

var keyMapping = map[int]ebiten.Key{
	vm.KeyShift:        ebiten.KeyShift,
	vm.KeyCtrl:         ebiten.KeyControl,
	vm.KeyAlt:          ebiten.KeyAlt,
	vm.KeyCap:          ebiten.KeyCapsLock,
	vm.KeyBack:         ebiten.KeyBackspace,
	vm.KeyTab:          ebiten.KeyTab,
	vm.KeyEnter:        ebiten.KeyEnter,
	vm.KeyF1:           ebiten.KeyF1,
	vm.KeyF2:           ebiten.KeyF2,
	vm.KeyF3:           ebiten.KeyF3,
	vm.KeyF4:           ebiten.KeyF4,
	vm.KeyF5:           ebiten.KeyF5,
	vm.KeyF6:           ebiten.KeyF6,
	vm.KeyF7:           ebiten.KeyF7,
	vm.KeyF8:           ebiten.KeyF8,
	vm.KeyF9:           ebiten.KeyF9,
	vm.KeyF10:          ebiten.KeyF10,
	vm.KeyF11:          ebiten.KeyF11,
	vm.KeyF12:          ebiten.KeyF12,
	vm.KeyLeft:         ebiten.KeyArrowLeft,
	vm.KeyRight:        ebiten.KeyArrowRight,
	vm.KeyUp:           ebiten.KeyArrowUp,
	vm.KeyDown:         ebiten.KeyArrowDown,
	vm.KeyEsc:          ebiten.KeyEscape,
	vm.KeySpace:        ebiten.KeySpace,
	vm.KeyApostrophe:   ebiten.KeyApostrophe,
	vm.KeyComma:        ebiten.KeyComma,
	vm.KeyMinus:        ebiten.KeyMinus,
	vm.KeyPeriod:       ebiten.KeyPeriod,
	vm.KeySlash:        ebiten.KeySlash,
	vm.KeyDigit0:       ebiten.KeyDigit0,
	vm.KeyDigit1:       ebiten.KeyDigit1,
	vm.KeyDigit2:       ebiten.KeyDigit2,
	vm.KeyDigit3:       ebiten.KeyDigit3,
	vm.KeyDigit4:       ebiten.KeyDigit4,
	vm.KeyDigit5:       ebiten.KeyDigit5,
	vm.KeyDigit6:       ebiten.KeyDigit6,
	vm.KeyDigit7:       ebiten.KeyDigit7,
	vm.KeyDigit8:       ebiten.KeyDigit8,
	vm.KeyDigit9:       ebiten.KeyDigit9,
	vm.KeySemicolon:    ebiten.KeySemicolon,
	vm.KeyEqual:        ebiten.KeyEqual,
	vm.KeyA:            ebiten.KeyA,
	vm.KeyB:            ebiten.KeyB,
	vm.KeyC:            ebiten.KeyC,
	vm.KeyD:            ebiten.KeyD,
	vm.KeyE:            ebiten.KeyE,
	vm.KeyF:            ebiten.KeyF,
	vm.KeyG:            ebiten.KeyG,
	vm.KeyH:            ebiten.KeyH,
	vm.KeyI:            ebiten.KeyI,
	vm.KeyJ:            ebiten.KeyJ,
	vm.KeyK:            ebiten.KeyK,
	vm.KeyL:            ebiten.KeyL,
	vm.KeyM:            ebiten.KeyM,
	vm.KeyN:            ebiten.KeyN,
	vm.KeyO:            ebiten.KeyO,
	vm.KeyP:            ebiten.KeyP,
	vm.KeyQ:            ebiten.KeyQ,
	vm.KeyR:            ebiten.KeyR,
	vm.KeyS:            ebiten.KeyS,
	vm.KeyT:            ebiten.KeyT,
	vm.KeyU:            ebiten.KeyU,
	vm.KeyV:            ebiten.KeyV,
	vm.KeyW:            ebiten.KeyW,
	vm.KeyX:            ebiten.KeyX,
	vm.KeyY:            ebiten.KeyY,
	vm.KeyZ:            ebiten.KeyZ,
	vm.KeyBracketLeft:  ebiten.KeyBracketLeft,
	vm.KeyBackslash:    ebiten.KeyBackslash,
	vm.KeyBracketRight: ebiten.KeyBracketRight,
	vm.KeyBackquote:    ebiten.KeyBackquote,
}

func updateKeyDuration() {
	for button, key := range keyMapping {
		if ebiten.IsKeyPressed(key) {
			vm.KeyDuration[button]++
		} else {
			vm.KeyDuration[button] = 0
		}
	}
}
