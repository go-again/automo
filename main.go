package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Point struct {
	X, Y int32
}

var (
	lastMousePosition Point
	lastMouseMoveTime time.Time
)

// Platform-specific interface
type mousePlatform interface {
	getCursorPos() (Point, error)
	setCursorPos(x, y int32) error
	moveMouseRelative(x, y int32) error
}

var mousePlatformImpl mousePlatform

func init() {
	mousePlatformImpl = newMouse()
}

func jiggleMouse(zenMode bool, diameter int32) {
	if zenMode {
		// Simulate mouse movement without actually moving cursor
		mousePlatformImpl.moveMouseRelative(0, 0)
		fmt.Println("Zen jiggle")
		return
	}

	startPos, _ := mousePlatformImpl.getCursorPos()
	radius := float64(diameter) / 2
	startTime := time.Now()
	duration := 200 * time.Millisecond

	// Number of steps to complete the circle
	steps := int(diameter)

	for i := 0; i <= steps; i++ {
		// Calculate progress (0 to 1)
		progress := float64(i) / float64(steps)

		// Calculate angle (0 to 2π)
		angle := progress * 2 * math.Pi

		// Calculate offset from center
		dx := int32(radius * math.Cos(angle))
		dy := int32(radius * math.Sin(angle))

		// Move to new position
		newX := startPos.X + dx
		newY := startPos.Y + dy
		mousePlatformImpl.setCursorPos(newX, newY)

		// Calculate remaining time and sleep accordingly
		elapsed := time.Since(startTime)
		if elapsed < duration {
			stepDuration := duration / time.Duration(steps)
			time.Sleep(stepDuration)
		}
	}

	// Return to starting position and prevent sleep
	mousePlatformImpl.setCursorPos(startPos.X, startPos.Y)
	mousePlatformImpl.moveMouseRelative(0, 0) // This will trigger the prevent sleep
	fmt.Println("Circle jiggle")
}

func checkMouseActivity(zenMode bool) {
	currentPos, _ := mousePlatformImpl.getCursorPos()

	if currentPos.X == lastMousePosition.X && currentPos.Y == lastMousePosition.Y {
		if time.Since(lastMouseMoveTime) >= 30*time.Second {
			jiggleMouse(zenMode, 5)
			lastMouseMoveTime = time.Now()
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
	currentPos, _ := mousePlatformImpl.getCursorPos()
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
