//go:build windows

package input

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Mouse events: https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-mouse_event

const (
	MOUSEEVENTF_MOVE     = 0x0001
	MOUSEEVENTF_ABSOLUTE = 0x8000
)

var (
	user32       = windows.NewLazyDLL("user32.dll")
	getCursorPos = user32.NewProc("GetCursorPos")
	mouseEvent   = user32.NewProc("mouse_event")

	lastPosition Point
)

type windowsPlatform struct {
	started bool
}

func NewPlatform() Platform {
	p := &windowsPlatform{}
	// Get initial cursor position
	pos, _ := p.GetCursorPos()
	lastPosition = pos
	p.started = true
	return p
}

func (p *windowsPlatform) GetCursorPos() (Point, error) {
	var point Point
	_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	return point, nil
}

func (p *windowsPlatform) SetCursorPos(pos Point) error {
	_, _, _ = mouseEvent.Call(MOUSEEVENTF_ABSOLUTE, uintptr(pos.X), uintptr(pos.Y), 0, 0)
	return nil
}

func (p *windowsPlatform) MoveCursorRelative(delta Point) error {
	_, _, _ = mouseEvent.Call(MOUSEEVENTF_MOVE, uintptr(delta.X), uintptr(delta.Y), 0, 0)
	return nil
}

func (p *windowsPlatform) HasUserActivity() bool {
	currentPos, _ := p.GetCursorPos()
	moved := currentPos.X != lastPosition.X || currentPos.Y != lastPosition.Y

	if Debug && moved {
		fmt.Printf("Mouse movement detected: %v -> %v\n", lastPosition, currentPos)
	}

	lastPosition = currentPos
	return moved
}
