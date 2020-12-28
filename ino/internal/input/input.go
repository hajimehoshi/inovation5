package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// TODO: This is duplicated with draw package's definitions.
	ScreenWidth  = 320
	ScreenHeight = 240
)

var theInput = &Input{
	pressed:     map[ebiten.Key]struct{}{},
	prevPressed: map[ebiten.Key]struct{}{},
}

func Current() *Input {
	return theInput
}

type Direction int

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionDown
)

var keys = []ebiten.Key{
	ebiten.KeyEnter,
	ebiten.KeySpace,
	ebiten.KeyLeft,
	ebiten.KeyDown,
	ebiten.KeyRight,

	// Fullscreen
	ebiten.KeyF,

	// Profiling
	ebiten.KeyP,
	ebiten.KeyQ,
}

type Input struct {
	pressed      map[ebiten.Key]struct{}
	prevPressed  map[ebiten.Key]struct{}
	touchEnabled bool

	gamepadID      ebiten.GamepadID
	gamepadEnabled bool
}

func (i *Input) IsTouchEnabled() bool {
	if isTouchPrimaryInput() {
		return true
	}
	return i.touchEnabled
}

func (i *Input) Update() {
	i.prevPressed = map[ebiten.Key]struct{}{}
	for k := range i.pressed {
		i.prevPressed[k] = struct{}{}
	}
	i.pressed = map[ebiten.Key]struct{}{}
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			i.pressed[k] = struct{}{}
		}
	}

	// Emulates the keys by gamepad pressing
	gamepadIDs := ebiten.GamepadIDs()
	if len(gamepadIDs) > 0 {
		if !i.gamepadEnabled {
			i.gamepadID = gamepadIDs[0]
			i.gamepadEnabled = true
		} else {
			var found bool
			for _, id := range gamepadIDs {
				if i.gamepadID == id {
					found = true
					break
				}
			}
			if !found {
				i.gamepadID = gamepadIDs[0]
			}
		}
	} else {
		i.gamepadEnabled = false
	}
	if i.gamepadEnabled {
		switch v := ebiten.GamepadAxis(i.gamepadID, 0); true {
		case -0.9 >= v:
			i.pressed[ebiten.KeyLeft] = struct{}{}
		case 0.9 <= v:
			i.pressed[ebiten.KeyRight] = struct{}{}
		}
		if y := ebiten.GamepadAxis(i.gamepadID, 1); y >= 0.9 {
			i.pressed[ebiten.KeyDown] = struct{}{}
		}

		for b := ebiten.GamepadButton0; b <= ebiten.GamepadButton3; b++ {
			if ebiten.IsGamepadButtonPressed(i.gamepadID, b) {
				i.pressed[ebiten.KeyEnter] = struct{}{}
				i.pressed[ebiten.KeySpace] = struct{}{}
				break
			}
		}
	}

	// Emulates the keys by touching
	touches := ebiten.TouchIDs()
	for _, t := range touches {
		x, y := ebiten.TouchPosition(t)
		// TODO(hajimehoshi): 64 are magic numbers
		if y < ScreenHeight-64 {
			continue
		}
		switch {
		case ScreenWidth <= x:
		case ScreenWidth*3/4 <= x:
			i.pressed[ebiten.KeyEnter] = struct{}{}
			i.pressed[ebiten.KeySpace] = struct{}{}
		case ScreenWidth*2/4 <= x:
			i.pressed[ebiten.KeyDown] = struct{}{}
		case ScreenWidth/4 <= x:
			i.pressed[ebiten.KeyRight] = struct{}{}
		default:
			i.pressed[ebiten.KeyLeft] = struct{}{}
		}
	}
	if 0 < len(touches) {
		i.touchEnabled = true
	}
}

func inLanguageSwitcher(x, y int) bool {
	return ScreenWidth*3/4 <= x && y < ScreenHeight/4
}

func (i *Input) IsSpaceTouched() bool {
	for _, t := range ebiten.TouchIDs() {
		x, y := ebiten.TouchPosition(t)
		if !inLanguageSwitcher(x, y) && y < ScreenHeight-64 {
			return true
		}
	}
	return false
}

func (i *Input) IsSpaceJustTouched() bool {
	for _, t := range inpututil.JustPressedTouchIDs() {
		x, y := ebiten.TouchPosition(t)
		if !inLanguageSwitcher(x, y) && y < ScreenHeight-64 {
			return true
		}
	}
	return false
}

func (i *Input) IsKeyPressed(key ebiten.Key) bool {
	_, ok := i.pressed[key]
	return ok
}

func (i *Input) IsKeyJustPressed(key ebiten.Key) bool {
	_, ok := i.pressed[key]
	if !ok {
		return false
	}
	_, ok = i.prevPressed[key]
	return !ok
}

func (i *Input) IsActionKeyPressed() bool {
	return i.IsKeyPressed(ebiten.KeyEnter) || i.IsKeyPressed(ebiten.KeySpace)
}

func (i *Input) IsActionKeyJustPressed() bool {
	return i.IsKeyJustPressed(ebiten.KeyEnter) || i.IsKeyJustPressed(ebiten.KeySpace)
}

func (i *Input) IsDirectionKeyPressed(dir Direction) bool {
	switch dir {
	case DirectionLeft:
		return i.IsKeyPressed(ebiten.KeyLeft)
	case DirectionRight:
		return i.IsKeyPressed(ebiten.KeyRight)
	case DirectionDown:
		return i.IsKeyPressed(ebiten.KeyDown)
	default:
		panic("not reach")
	}
}

func (i *Input) IsLanguageSwitcherPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if inLanguageSwitcher(ebiten.CursorPosition()) {
			return true
		}
	}
	for _, t := range inpututil.JustPressedTouchIDs() {
		if inLanguageSwitcher(ebiten.TouchPosition(t)) {
			return true
		}
	}
	return false
}
