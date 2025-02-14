package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Point struct {
	X, Y int32
}

var (
	user32            = windows.NewLazyDLL("user32.dll")
	getCursorPos      = user32.NewProc("GetCursorPos")
	setCursorPos      = user32.NewProc("SetCursorPos")
	mouse_event       = user32.NewProc("mouse_event")
	lastMousePosition Point
	lastMouseMoveTime time.Time
	zigzag            bool // Alternates between true/false for zigzag movement
)

const (
	MOUSEEVENTF_MOVE = 0x0001
)

func jiggleMouse(zenMode bool, distance int32) {
	if zenMode {
		// Simulate mouse movement without actually moving cursor
		_, _, _ = mouse_event.Call(MOUSEEVENTF_MOVE, 0, 0, 0, 0)
		fmt.Println("Zen jiggle")
	} else {
		var currentPos Point
		_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(&currentPos)))

		// Zigzag movement pattern
		movement := distance
		if !zigzag {
			movement = -distance
		}
		zigzag = !zigzag

		newX := currentPos.X + movement
		newY := currentPos.Y + movement
		_, _, _ = setCursorPos.Call(uintptr(newX), uintptr(newY))
		fmt.Printf("Mouse moved to: %d, %d (zigzag: %v)\n", newX, newY, zigzag)
	}
	lastMouseMoveTime = time.Now()
}

func checkMouseActivity(zenMode bool) {
	var currentPos Point
	_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(&currentPos)))

	if currentPos.X == lastMousePosition.X && currentPos.Y == lastMousePosition.Y {
		if time.Since(lastMouseMoveTime) >= 30*time.Second {
			jiggleMouse(zenMode, 4)
		}
	} else {
		lastMousePosition = currentPos
		lastMouseMoveTime = time.Now()
	}
}

func main() {
	// Command line flags
	zenMode := flag.Bool("zen", false, "Enable zen mode (virtual mouse movement)")
	checkInterval := flag.Int("interval", 5, "Check interval in seconds")
	flag.Parse()

	fmt.Printf("Mouse mover started in %s mode. Press Ctrl+C to exit.\n",
		map[bool]string{true: "zen", false: "normal"}[*zenMode])
	fmt.Printf("Check interval: %d seconds\n", *checkInterval)

	// Create a channel to handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for periodic checks
	ticker := time.NewTicker(time.Duration(*checkInterval) * time.Second)
	defer ticker.Stop()

	// Initialize the last mouse position and time
	var currentPos Point
	_, _, _ = getCursorPos.Call(uintptr(unsafe.Pointer(&currentPos)))
	lastMousePosition = currentPos
	lastMouseMoveTime = time.Now()

	// Main loop
	for {
		select {
		case <-ticker.C:
			checkMouseActivity(*zenMode)
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		}
	}
}
