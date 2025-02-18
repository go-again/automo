//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ApplicationServices
#include <CoreGraphics/CoreGraphics.h>

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

void preventSleep() {
    // Create and post a minimal middle mouse down and up event
	// https://developer.apple.com/documentation/coregraphics/cgmousebutton?changes=_9&language=objc
    CGPoint currentPos = getMousePos();
    //CGEventRef mouseDown = CGEventCreateMouseEvent(NULL, kCGEventOtherMouseDown, currentPos, kCGMouseButtonCenter);
    //CGEventRef mouseUp = CGEventCreateMouseEvent(NULL, kCGEventOtherMouseUp, currentPos, kCGMouseButtonCenter);
	CGEventRef mouseDown = CGEventCreateMouseEvent(NULL, kCGEventOtherMouseDown, currentPos, 31);
    CGEventRef mouseUp = CGEventCreateMouseEvent(NULL, kCGEventOtherMouseUp, currentPos, 31);

    // Post the events with a tiny delay
    CGEventPost(kCGHIDEventTap, mouseDown);
    usleep(100000); // 100ms delay
    CGEventPost(kCGHIDEventTap, mouseUp);

    // Clean up
    CFRelease(mouseDown);
    CFRelease(mouseUp);
}
*/
import "C"

type darwinMouse struct{}

func newMouse() mousePlatform {
	return &darwinMouse{}
}

func getMouseLocation() (x, y float64) {
	point := C.getMousePos()
	return float64(point.x), float64(point.y)
}

func moveMouseAbsolute(x, y float64) {
	C.moveMouse(C.CGFloat(x), C.CGFloat(y))
}

func moveMouseRelative(x, y float64) {
	currentX, currentY := getMouseLocation()
	moveMouseAbsolute(currentX+x, currentY+y)
}

func preventScreenLock() {
	C.preventSleep()
}

func (d *darwinMouse) getCursorPos() (Point, error) {
	x, y := getMouseLocation()
	return Point{X: int32(x), Y: int32(y)}, nil
}

func (d *darwinMouse) setCursorPos(x, y int32) error {
	moveMouseAbsolute(float64(x), float64(y))
	return nil
}

func (d *darwinMouse) moveMouseRelative(x, y int32) error {
	if x == 0 && y == 0 {
		// For zen mode, simulate activity without visible movement
		preventScreenLock()
	} else {
		moveMouseRelative(float64(x), float64(y))
		// Also prevent sleep after movement
		preventScreenLock()
	}
	return nil
}
