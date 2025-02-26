package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"automo/input"
)

const (
	// ActivityTimeout is the time after which the mouse will jiggle if there is no user activity
	ActivityTimeout = 30 * time.Second
	// JiggleDuration is the duration of the mouse jiggle
	JiggleDuration = 200 * time.Millisecond
	// JiggleDiameter is the diameter of the circle that the mouse will jiggle in
	JiggleDiameter = 5
)

var lastActivityTime time.Time

// parseFlags parses command-line flags
func parseFlags() (bool, int, bool) {
	zenMode := flag.Bool("zen", false, "Enable zen mode (virtual mouse movement)")
	checkInterval := flag.Int("interval", 5, "Check interval in seconds")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()
	return *zenMode, *checkInterval, *debug
}

// setupPlatform sets up the platform
func setupPlatform() (input.Platform, error) {
	return input.New()
}

// jiggleMouse jiggles the mouse
func jiggleMouse(zenMode bool, diameter int32, platform input.Platform) error {
	if zenMode {
		// Simulate mouse movement without actually moving the cursor
		return platform.MoveCursorRelative(input.Point{0, 0})
	}

	startPos, err := platform.GetCursorPos()
	if err != nil {
		return err
	}
	radius := float64(diameter) / 2
	startTime := time.Now()
	duration := JiggleDuration

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
		newPos := input.Point{
			X: startPos.X + dx,
			Y: startPos.Y + dy,
		}
		err = platform.SetCursorPos(newPos)
		if err != nil {
			return err
		}

		// Calculate remaining time and sleep accordingly
		elapsed := time.Since(startTime)
		if elapsed < duration {
			stepDuration := duration / time.Duration(steps)
			time.Sleep(stepDuration)
		}
	}

	// Return to starting position and prevent sleep
	err = platform.SetCursorPos(startPos)
	if err != nil {
		return err
	}
	return platform.MoveCursorRelative(input.Point{0, 0}) // This will trigger the prevent sleep
}

// checkMouseActivity checks for mouse activity
func checkMouseActivity(zenMode bool, platform input.Platform) error {
	if platform.HasUserActivity() {
		lastActivityTime = time.Now()
	} else {
		if time.Since(lastActivityTime) >= ActivityTimeout {
			err := jiggleMouse(zenMode, JiggleDiameter, platform)
			if err != nil {
				return err
			}
			lastActivityTime = time.Now()
		}
	}
	return nil
}

// runMainLoop runs the main loop
func runMainLoop(zenMode bool, checkInterval int, platform input.Platform) {
	// Create a channel to handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for periodic checks
	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	defer ticker.Stop()

	lastActivityTime = time.Now()

	// Main loop
	for {
		select {
		case <-ticker.C:
			err := checkMouseActivity(zenMode, platform)
			if err != nil {
				fmt.Println(err)
			}
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		}
	}
}

func main() {
	zenMode, checkInterval, debug := parseFlags()
	input.Debug = debug

	platform, err := setupPlatform()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("automo started in %s mode. Press Ctrl+C to exit.\n",
		map[bool]string{true: "zen", false: "normal"}[zenMode])
	fmt.Printf("Check interval: %d seconds\n", checkInterval)
	if debug {
		fmt.Println("Debug output enabled")
	}

	defer func() {
		if err := platform.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	runMainLoop(zenMode, checkInterval, platform)
}
