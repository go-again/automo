//go:build windows

package input

import (
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MOUSEEVENTF_MOVE = 0x0001
	WH_KEYBOARD_LL   = 13
	WH_MOUSE_LL      = 14
	WM_KEYDOWN       = 0x0100
	WM_KEYUP         = 0x0101
	WM_MOUSEMOVE     = 0x0200
	WM_LBUTTONDOWN   = 0x0201
	WM_LBUTTONUP     = 0x0202
	WM_RBUTTONDOWN   = 0x0204
	WM_RBUTTONUP     = 0x0205
	WM_MOUSEWHEEL    = 0x020A
)

var (
	user32       = windows.NewLazyDLL("user32.dll")
	getCursorPos = user32.NewProc("GetCursorPos")
	setCursorPos = user32.NewProc("SetCursorPos")
	mouse_event  = user32.NewProc("mouse_event")
	setHook      = user32.NewProc("SetWindowsHookExW")
	getMsg       = user32.NewProc("GetMessageW")
	unhook       = user32.NewProc("UnhookWindowsHookEx")

	keyboardHook   windows.Handle
	mouseHook      windows.Handle
	userActivity   uint32
	hookThreadDone = make(chan bool)

	translateMessage = user32.NewProc("TranslateMessage")
	dispatchMessage  = user32.NewProc("DispatchMessageW")
	callNextHookEx   = user32.NewProc("CallNextHookEx")
)

type MSG struct {
	Hwnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      Point
}

type windowsPlatform struct {
	started bool
}

func NewPlatform() Platform {
	p := &windowsPlatform{}
	if err := startHooks(); err != nil {
		panic("Failed to start input hooks: " + err.Error())
	}
	p.started = true
	runtime.SetFinalizer(p, func(p *windowsPlatform) {
		if p.started {
			stopHooks()
		}
	})
	return p
}

func startHooks() error {
	var err error
	keyboardHook, err = setLowLevelHook(WH_KEYBOARD_LL)
	if err != nil {
		return fmt.Errorf("failed to set keyboard hook: %v", err)
	}

	mouseHook, err = setLowLevelHook(WH_MOUSE_LL)
	if err != nil {
		unhook.Call(uintptr(keyboardHook))
		return fmt.Errorf("failed to set mouse hook: %v", err)
	}

	// Start message loop in a separate goroutine
	go messageLoop()
	return nil
}

func stopHooks() {
	if keyboardHook != 0 {
		unhook.Call(uintptr(keyboardHook))
		keyboardHook = 0
	}
	if mouseHook != 0 {
		unhook.Call(uintptr(mouseHook))
		mouseHook = 0
	}
}

func setLowLevelHook(hookType int) (windows.Handle, error) {
	hook, _, err := setHook.Call(
		uintptr(hookType),
		windows.NewCallback(hookCallback),
		0,
		0,
	)
	if hook == 0 {
		return 0, err
	}
	return windows.Handle(hook), nil
}

func hookCallback(code int, wparam, lparam uintptr) uintptr {
	if code >= 0 {
		switch wparam {
		case WM_KEYDOWN, WM_KEYUP:
			userActivity |= uint32(ActivityKeyboard)
		case WM_MOUSEMOVE:
			userActivity |= uint32(ActivityMouseMove)
		case WM_LBUTTONDOWN, WM_LBUTTONUP, WM_RBUTTONDOWN, WM_RBUTTONUP:
			userActivity |= uint32(ActivityMouseClick)
		case WM_MOUSEWHEEL:
			userActivity |= uint32(ActivityScroll)
		}
	}
	// Use the imported callNextHookEx proc instead
	ret, _, _ := callNextHookEx.Call(0, uintptr(code), wparam, lparam)
	return ret
}

func messageLoop() {
	var msg MSG
	for {
		ret, _, _ := getMsg.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
	hookThreadDone <- true
}

// Rename receiver to match new type name
func (p *windowsPlatform) GetCursorPos() (Point, error) {
	var point Point
	_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(&point)))
	return point, nil
}

func (p *windowsPlatform) SetCursorPos(pos Point) error {
	_, _, _ = setCursorPos.Call(uintptr(pos.X), uintptr(pos.Y))
	return nil
}

func (p *windowsPlatform) MoveCursorRelative(delta Point) error {
	_, _, _ = mouse_event.Call(MOUSEEVENTF_MOVE, uintptr(delta.X), uintptr(delta.Y), 0, 0)
	return nil
}

func (p *windowsPlatform) HasUserActivity() bool {
	activity := ActivityType(userActivity)
	if Debug && (activity != ActivityNone) {
		fmt.Printf("Activity detected: %d\n", activity)
	}
	userActivity = 0
	return activity != ActivityNone
}
