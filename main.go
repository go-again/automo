package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"automo/input"
)

var (
	lastActivityTime time.Time
)

func jiggleMouse(zenMode bool, diameter int32, platform input.Platform) {
	if zenMode {
		// Simulate mouse movement without actually moving cursor
		platform.MoveCursorRelative(input.Point{0, 0})
		fmt.Println("Zen jiggle")
		return
	}

	startPos, _ := platform.GetCursorPos()
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
		newPos := input.Point{
			X: startPos.X + dx,
			Y: startPos.Y + dy,
		}
		platform.SetCursorPos(newPos)

		// Calculate remaining time and sleep accordingly
		elapsed := time.Since(startTime)
		if elapsed < duration {
			stepDuration := duration / time.Duration(steps)
			time.Sleep(stepDuration)
		}
	}

	// Return to starting position and prevent sleep
	platform.SetCursorPos(startPos)
	platform.MoveCursorRelative(input.Point{0, 0}) // This will trigger the prevent sleep
	fmt.Println("Circle jiggle")
}

func checkMouseActivity(zenMode bool, platform input.Platform) {
	if platform.HasUserActivity() {
		lastActivityTime = time.Now()
	} else {
		if time.Since(lastActivityTime) >= 30*time.Second {
			jiggleMouse(zenMode, 5, platform)
			lastActivityTime = time.Now()
		}
	}
}

func main() {
	// Command line flags
	zenMode := flag.Bool("zen", false, "Enable zen mode (virtual mouse movement)")
	checkInterval := flag.Int("interval", 5, "Check interval in seconds")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	// Set debug mode
	input.Debug = *debug

	platform := input.New()

	fmt.Printf("automo started in %s mode. Press Ctrl+C to exit.\n",
		map[bool]string{true: "zen", false: "normal"}[*zenMode])
	fmt.Printf("Check interval: %d seconds\n", *checkInterval)
	if *debug {
		fmt.Println("Debug output enabled")
	}

	// Create a channel to handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for periodic checks
	ticker := time.NewTicker(time.Duration(*checkInterval) * time.Second)
	defer ticker.Stop()

	lastActivityTime = time.Now()

	// Main loop
	for {
		select {
		case <-ticker.C:
			checkMouseActivity(*zenMode, platform)
		case <-sigChan:
			fmt.Println("\nShutting down...")
			return
		}
	}
}
