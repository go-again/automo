//go:build darwin

package input

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ApplicationServices
#include <CoreGraphics/CoreGraphics.h>
#include <stdint.h>
#include <stdio.h>
#include <pthread.h>
#include <unistd.h>

static CFMachPortRef eventTap;
static CFRunLoopSourceRef runLoopSource;
static volatile uint32_t userActivity = 0;

CGEventRef eventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    switch (type) {
        case kCGEventNull:
            break;

        case kCGEventKeyDown:
        case kCGEventKeyUp:
        case kCGEventFlagsChanged:
            userActivity |= 1;
            break;

        case kCGEventMouseMoved:
        case kCGEventLeftMouseDragged:
        case kCGEventRightMouseDragged:
        case kCGEventOtherMouseDragged:
            userActivity |= 2;
            break;

        case kCGEventLeftMouseDown:
        case kCGEventLeftMouseUp:
        case kCGEventRightMouseDown:
        case kCGEventRightMouseUp:
        case kCGEventOtherMouseDown:
        case kCGEventOtherMouseUp:
            userActivity |= 4;
            break;

        case kCGEventScrollWheel:
            userActivity |= 8;
            break;

        case kCGEventTabletPointer:
        case kCGEventTabletProximity:
        case kCGEventTapDisabledByTimeout:
        case kCGEventTapDisabledByUserInput:
            break;
    }
    return event;
}

static CFRunLoopRef tapLoop = NULL;
static volatile bool shouldRun = true;

void* runEventLoop(void* arg) {
    CGEventMask eventMask = (
        CGEventMaskBit(kCGEventKeyDown) |
        CGEventMaskBit(kCGEventKeyUp) |
        CGEventMaskBit(kCGEventFlagsChanged) |
        CGEventMaskBit(kCGEventMouseMoved) |
        CGEventMaskBit(kCGEventLeftMouseDragged) |
        CGEventMaskBit(kCGEventRightMouseDragged) |
        CGEventMaskBit(kCGEventOtherMouseDragged) |
        CGEventMaskBit(kCGEventLeftMouseDown) |
        CGEventMaskBit(kCGEventLeftMouseUp) |
        CGEventMaskBit(kCGEventRightMouseDown) |
        CGEventMaskBit(kCGEventRightMouseUp) |
        CGEventMaskBit(kCGEventOtherMouseDown) |
        CGEventMaskBit(kCGEventOtherMouseUp) |
        CGEventMaskBit(kCGEventScrollWheel)
    );

    eventTap = CGEventTapCreate(
        kCGSessionEventTap,
        kCGHeadInsertEventTap,
        kCGEventTapOptionListenOnly,
        eventMask,
        eventCallback,
        NULL
    );

    if (!eventTap) {
        printf("Failed to create event tap. Make sure the app has accessibility permissions.\n");
        return NULL;
    }

    runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
    if (!runLoopSource) {
        printf("Failed to create run loop source\n");
        CFRelease(eventTap);
        return NULL;
    }

    tapLoop = CFRunLoopGetCurrent();
    CFRunLoopAddSource(tapLoop, runLoopSource, kCFRunLoopCommonModes);
    CGEventTapEnable(eventTap, true);

    while (shouldRun) {
        CFRunLoopRunInMode(kCFRunLoopDefaultMode, 1.0, false);
    }

    return NULL;
}

bool startEventTap() {
    pthread_t thread;
    shouldRun = true;
    if (pthread_create(&thread, NULL, runEventLoop, NULL) != 0) {
        printf("Failed to create event monitoring thread\n");
        return false;
    }
    pthread_detach(thread);
    return true;
}

void stopEventTap() {
    shouldRun = false;

    if (eventTap) {
        CGEventTapEnable(eventTap, false);
    }

    if (runLoopSource) {
        if (tapLoop) {
            CFRunLoopRemoveSource(tapLoop, runLoopSource, kCFRunLoopCommonModes);
        }
        CFRelease(runLoopSource);
        runLoopSource = NULL;
    }

    if (eventTap) {
        CFRelease(eventTap);
        eventTap = NULL;
    }

    tapLoop = NULL;
}

uint32_t checkAndClearActivity() {
    uint32_t current = userActivity;
    userActivity = 0;
    return current;
}

CGPoint getMousePos() {
    CGEventRef event = CGEventCreate(NULL);
    CGPoint point = CGEventGetLocation(event);
    CFRelease(event);
    return point;
}

void moveMouse(CGFloat x, CGFloat y) {
    CGPoint point;
    point.x = x;
    point.y = y;
    CGWarpMouseCursorPosition(point);
    CGAssociateMouseAndMouseCursorPosition(true);
}

bool hasKeyboardActivity() {
    // Get the current keyboard state
    CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
    if (!source) return false;

    // Check a range of common keys
    bool activity = false;
    for (int i = 0; i < 128 && !activity; i++) {
        if (CGEventSourceKeyState(kCGEventSourceStateHIDSystemState, i)) {
            activity = true;
            break;
        }
    }

    if (source) {
        CFRelease(source);
    }
    return activity;
}

void simulateF20Press() {
    CGEventRef keyDown = CGEventCreateKeyboardEvent(NULL, 0x6D, true);  // 0x6D is F20
    CGEventRef keyUp = CGEventCreateKeyboardEvent(NULL, 0x6D, false);

    CGEventPost(kCGHIDEventTap, keyDown);
    usleep(50000); // 50ms delay
    CGEventPost(kCGHIDEventTap, keyUp);

    CFRelease(keyDown);
    CFRelease(keyUp);
}

void preventSleep() {
    // Simulate F20 key press to prevent sleep
    simulateF20Press();
}
*/
import "C"
import (
	"fmt"
	"runtime"
)

type darwinPlatform struct {
	started bool
}

func NewPlatform() Platform {
	p := &darwinPlatform{}
	if ok := C.startEventTap(); !ok {
		panic("Failed to start event tap")
	}
	p.started = true
	runtime.SetFinalizer(p, func(p *darwinPlatform) {
		if p.started {
			C.stopEventTap()
		}
	})
	return p
}

// Helper functions with consistent naming
func getCursorLocation() (x, y float64) {
	point := C.getMousePos()
	return float64(point.x), float64(point.y)
}

func moveCursorAbsolute(x, y float64) {
	C.moveMouse(C.CGFloat(x), C.CGFloat(y))
}

func moveCursorRelative(x, y float64) {
	currentX, currentY := getCursorLocation()
	moveCursorAbsolute(currentX+x, currentY+y)
}

func preventScreenLock() {
	C.preventSleep()
}

func (p *darwinPlatform) GetCursorPos() (Point, error) {
	x, y := getCursorLocation()
	return Point{X: int32(x), Y: int32(y)}, nil
}

func (p *darwinPlatform) SetCursorPos(pos Point) error {
	moveCursorAbsolute(float64(pos.X), float64(pos.Y))
	return nil
}

func (p *darwinPlatform) MoveCursorRelative(delta Point) error {
	if delta.X == 0 && delta.Y == 0 {
		// For zen mode, simulate activity without visible movement
		preventScreenLock()
	} else {
		moveCursorRelative(float64(delta.X), float64(delta.Y))
		preventScreenLock()
	}
	return nil
}

func (p *darwinPlatform) HasUserActivity() bool {
	activity := ActivityType(C.checkAndClearActivity())
	if Debug && (activity != ActivityNone) {
		fmt.Printf("Activity detected: %d\n", activity)
	}
	return activity != ActivityNone
}

// Optional: Add method to get detailed activity info if needed
func (d *darwinPlatform) GetActivityDetails() ActivityType {
	return ActivityType(C.checkAndClearActivity())
}
