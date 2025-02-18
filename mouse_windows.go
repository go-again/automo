//go:build windows

package main

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	MOUSEEVENTF_MOVE = 0x0001
)

var (
	user32       = windows.NewLazyDLL("user32.dll")
	getCursorPos = user32.NewProc("GetCursorPos")
	setCursorPos = user32.NewProc("SetCursorPos")
	mouse_event  = user32.NewProc("mouse_event")
)

type windowsMouse struct{}

func newMouse() mousePlatform {
	return &windowsMouse{}
}

func (w *windowsMouse) getCursorPos() (Point, error) {
	var point Point
	_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	return point, nil
}

func (w *windowsMouse) setCursorPos(x, y int32) error {
	_, _, _ = setCursorPos.Call(uintptr(x), uintptr(y))
	return nil
}

func (w *windowsMouse) moveMouseRelative(x, y int32) error {
	_, _, _ = mouse_event.Call(MOUSEEVENTF_MOVE, uintptr(x), uintptr(y), 0, 0)
	return nil
}
