//go:build windows

package input

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

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

// NewPlatform returns a new platform-specific implementation
func NewPlatform() (Platform, error) {
	p := &windowsPlatform{}
	// Get initial cursor position
	pos, err := p.GetCursorPos()
	if err != nil {
		return nil, err
	}
	lastPosition = pos
	p.started = true
	return p, nil
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
	currentPos, err := p.GetCursorPos()
	if err != nil {
		return false
	}
	moved := currentPos.X != lastPosition.X || currentPos.Y != lastPosition.Y

	if Debug && moved {
		fmt.Printf("Mouse movement detected: %v -> %v\n", lastPosition, currentPos)
	}

	lastPosition = currentPos
	return moved
}

func (p *windowsPlatform) Close() error {
	// No resources to release on Windows
	return nil
}
