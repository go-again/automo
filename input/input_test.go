package input

import (
	"fmt"
	"testing"
)

// MockPlatform implements the Platform interface for testing purposes
type MockPlatform struct {
	cursorX  int32
	cursorY  int32
	activity bool
}

// NewMockPlatform returns a new instance of MockPlatform
func NewMockPlatform() *MockPlatform {
	return &MockPlatform{
		cursorX:  0,
		cursorY:  0,
		activity: false,
	}
}

// GetCursorPos returns the current cursor position
func (m *MockPlatform) GetCursorPos() (Point, error) {
	return Point{X: m.cursorX, Y: m.cursorY}, nil
}

// SetCursorPos sets the cursor position
func (m *MockPlatform) SetCursorPos(pos Point) error {
	m.cursorX = pos.X
	m.cursorY = pos.Y
	return nil
}

// MoveCursorRelative moves the cursor relative to its current position
func (m *MockPlatform) MoveCursorRelative(delta Point) error {
	m.cursorX += delta.X
	m.cursorY += delta.Y
	return nil
}

// HasUserActivity returns whether there is user activity
func (m *MockPlatform) HasUserActivity() bool {
	return m.activity
}

func TestNew(t *testing.T) {
	// Test that New returns a non-nil platform implementation
	platform, err := New()
	if err != nil {
		t.Error("New returned an error:", err)
	}
	if platform == nil {
		t.Error("New returned nil platform implementation")
	}
}

func TestPoint(t *testing.T) {
	// Test that Point correctly holds X and Y values
	point := Point{X: 10, Y: 20}
	if point.X != 10 || point.Y != 20 {
		t.Error("Point does not hold correct X and Y values")
	}
}

func TestActivityType(t *testing.T) {
	// Test that ActivityType constants have the correct values
	tests := []struct {
		name         string
		activityType ActivityType
		expected     uint32
	}{
		{"ActivityNone", ActivityNone, 0},
		{"ActivityKeyboard", ActivityKeyboard, 1},
		{"ActivityMouseMove", ActivityMouseMove, 2},
		{"ActivityMouseClick", ActivityMouseClick, 4},
		{"ActivityScroll", ActivityScroll, 8},
	}

	for _, test := range tests {
		if uint32(test.activityType) != test.expected {
			t.Errorf("%s should be %d, but got %d", test.name, test.expected, uint32(test.activityType))
		}
	}
}

func TestPlatformMethods(t *testing.T) {
	platform := NewMockPlatform()

	// Test GetCursorPos
	initialPos, err := platform.GetCursorPos()
	if err != nil {
		t.Error("GetCursorPos returned an error")
	}
	if initialPos.X != 0 || initialPos.Y != 0 {
		t.Error("GetCursorPos did not return (0,0) as initial position")
	}

	// Test SetCursorPos
	newPos := Point{X: 100, Y: 200}
	err = platform.SetCursorPos(newPos)
	if err != nil {
		t.Error("SetCursorPos returned an error")
	}

	updatedPos, err := platform.GetCursorPos()
	if err != nil {
		t.Error("GetCursorPos returned an error after SetCursorPos")
	}
	if updatedPos.X != newPos.X || updatedPos.Y != newPos.Y {
		t.Error("SetCursorPos did not update the cursor position correctly")
	}

	// Test MoveCursorRelative
	delta := Point{X: 50, Y: -100}
	err = platform.MoveCursorRelative(delta)
	if err != nil {
		t.Error("MoveCursorRelative returned an error")
	}

	newPosAfterMove, err := platform.GetCursorPos()
	if err != nil {
		t.Error("GetCursorPos returned an error after MoveCursorRelative")
	}
	expectedX := updatedPos.X + delta.X
	expectedY := updatedPos.Y + delta.Y
	if newPosAfterMove.X != expectedX || newPosAfterMove.Y != expectedY {
		t.Error("MoveCursorRelative did not move the cursor correctly")
	}

	// Test HasUserActivity
	platform.activity = true
	if !platform.HasUserActivity() {
		t.Error("HasUserActivity returned false when activity is true")
	}
	platform.activity = false
	if platform.HasUserActivity() {
		t.Error("HasUserActivity returned true when activity is false")
	}
}

func TestMain(m *testing.M) {
	// Any setup code can go here
	fmt.Println("Starting tests...")
	m.Run()
	// Any cleanup code can go here
	fmt.Println("Tests completed.")
}
