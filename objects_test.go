package gologo

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

const epsilon = 0.00001

var createObjTests = []struct {
	name       string
	posX, posY float32
	zOrder     int
}{
	{"on screen", 200, 300, 1},
	{"at origin", 0, 0, 2},
	{"off screen", -100, -200, 3},
}

// TestCreateObject : Test basic Object creation by building the object
// then reading values back out to confirm values
// A simple object has no renderer
func TestCreateObject(t *testing.T) {
	var position mgl32.Vec3
	var obj *Object
	var x, y float32
	var age, direction, zOrder int

	for _, tc := range createObjTests {
		t.Run(tc.name, func(t *testing.T) {
			position = mgl32.Vec3{tc.posX, tc.posY, 0.0}
			obj = CreateObject(position)
			obj.SetZOrder(tc.zOrder)

			// Test default position
			x, y = obj.GetPosition()
			if x != tc.posX || y != tc.posY {
				t.Errorf("Location was x (%v), y (%v) should be x (%v), y (%v)",
					x, y, tc.posX, tc.posY)
			}

			// Test default direction
			direction = radToNearestDeg(obj.Direction())
			if direction%360 != 0 {
				t.Errorf("Direction was (%v) should be 0", direction)
			}

			// Test object has age
			age = obj.GetAge()
			if age < 0 {
				t.Errorf("Age was (%v) should not be negative", age)
			}

			// Test default zOrder
			zOrder = obj.GetZOrder()
			if zOrder != tc.zOrder {
				t.Errorf("zOrder was (%v) should be (%v)", zOrder, tc.zOrder)
			}
		})
	}
}

var translateObjTests = []struct {
	name                                         string
	posX, posY, transX, transY, posXExp, posYExp float32
}{
	{"up and right from origin", 0, 0, 100, 200, 100, 200},
	{"up and right not origin", 100, 150, 200, 200, 300, 350},
	{"down and left not origin", 100, 150, -200, -450, -100, -300},
	{"to origin", -100, -200, 100, 200, 0, 0},
}

// TestObjectTranslate : Test that Object translation moves the object by the specified amount
func TestObjectTranslate(t *testing.T) {
	var position mgl32.Vec3
	var x, y float32
	var obj *Object

	for _, tc := range translateObjTests {
		t.Run(tc.name, func(t *testing.T) {
			position = mgl32.Vec3{tc.posX, tc.posY, 0.0}
			obj = CreateObject(position)

			// Test translation
			obj.Translate(tc.transX, tc.transY)
			x, y = obj.GetPosition()
			if x != tc.posXExp || y != tc.posYExp {
				t.Errorf("After translation location was x (%v), y (%v) should be x (%v), y (%v)",
					x, y, tc.posXExp, tc.posYExp)
			}
		})
	}
}

var rotateObjTests = []struct {
	name                           string
	startAngle, rotation, angleExp float64
}{
	{"rotate 90 from 0", 0, math.Pi / 2, math.Pi / 2},
	{"rotate -90 from 0", 0, -math.Pi / 2, -math.Pi / 2},
	{"rotate 360 from 0", 0, math.Pi * 2, 0},
	{"rotate 45 from 90", math.Pi / 2, math.Pi / 4, 3 * math.Pi / 4},
	{"rotate 90 from 270", 3 * math.Pi / 2, math.Pi / 2, 0},
}

// TestObjectRotate : Test that Object rotation changes the orientation of the object by the specified amount
func TestObjectRotate(t *testing.T) {
	var direction, diff float64
	var obj *Object

	for _, tc := range rotateObjTests {
		t.Run(tc.name, func(t *testing.T) {
			obj = CreateObject(mgl32.Vec3{0.0, 0.0, 0.0})

			// initialise angle
			if tc.startAngle != 0 {
				obj.Rotate(tc.startAngle)
			}

			// Test rotating by 90 degrees
			obj.Rotate(tc.rotation)
			direction = obj.Direction()
			diff = math.Abs(direction - tc.angleExp)
			if diff > epsilon {
				t.Errorf("Direction was (%v) should be (%v) with tolerance (%v)", direction, tc.angleExp, epsilon)
			}
		})
	}
}
