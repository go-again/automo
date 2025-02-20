package input

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

var (
	Debug bool
)

// Platform represents platform-specific input device functionality
type Platform interface {
	GetCursorPos() (Point, error)
	SetCursorPos(pos Point) error
	MoveCursorRelative(delta Point) error
	HasUserActivity() bool
}

// New returns a new platform-specific implementation
func New() Platform {
	return NewPlatform()
}
