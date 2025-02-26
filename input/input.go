package input

import "errors"

// Point represents a 2D coordinate
type Point struct {
	X, Y int32
}

// ActivityType represents different types of user activity
type ActivityType uint32

const (
	ActivityNone       ActivityType = 0
	ActivityKeyboard   ActivityType = 1
	ActivityMouseMove  ActivityType = 2
	ActivityMouseClick ActivityType = 4
	ActivityScroll     ActivityType = 8
)

// Debug flag to enable debug logging
var Debug bool

// Platform represents platform-specific input device functionality
type Platform interface {
	// GetCursorPos returns the current cursor position
	GetCursorPos() (Point, error)
	// SetCursorPos sets the cursor position to the specified point
	SetCursorPos(pos Point) error
	// MoveCursorRelative moves the cursor by the specified delta
	MoveCursorRelative(delta Point) error
	// HasUserActivity checks if there has been any user activity
	HasUserActivity() bool
	// Close releases any resources held by the platform
	Close() error
}

// New returns a new platform-specific implementation
func New() (Platform, error) {
	return NewPlatform()
}

// ErrNotImplemented is returned when a function is not implemented
var ErrNotImplemented = errors.New("not implemented")
