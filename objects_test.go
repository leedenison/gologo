package gologo

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	positionX, positionY float32 = 500, 400
)

// TestSimpleObjectCreation : Test basic Object creation by building the object
// then reading values back out to confirm values
// A simple object has no renderer
func TestSimpleObjectCreation(t *testing.T) {
	position := mgl32.Vec3{positionX, positionY, 0.0}
	obj := CreateObject(position)
	obj.SetZOrder(1)

	// Test default position
	x, y := obj.GetPosition()
	if x != positionX || y != positionY {
		t.Errorf("Location was x (%v), y (%v) should be x (%v), y (%v)",
			x, y, positionX, positionY)
	}

	// Test default direction
	direction := radToNearestDeg(obj.Direction())
	if direction%360 != 0 {
		t.Errorf("Direction was (%v) should be 0", direction)
	}

	// Test object has age
	age := obj.GetAge()
	if age < 0 {
		t.Errorf("Age was (%v) should not be negative", age)
	}
}

// TestObjectTranslation : Test that Object translation moves the object by the specified amount
func TestObjectTranslation(t *testing.T) {
	var translateX, translateY float32 = 150, 100
	position := mgl32.Vec3{positionX, positionY, 0.0}
	obj := CreateObject(position)

	// Test translation
	obj.Translate(translateX, translateY)
	x, y := obj.GetPosition()
	if x != positionX+translateX || y != positionY+translateY {
		t.Errorf("After translation location was x (%v), y (%v) should be x (%v), y (%v)",
			x, y, positionX+translateX, positionY+translateY)
	}
}


// TestObjectRotation : Test that Object rotation changes the orientation of the object by the specified amount
func TestObjectRotation(t *testing.T) {
	position := mgl32.Vec3{positionX, positionY, 0.0}
	obj := CreateObject(position)

	// Test rotating by 90 degrees
	obj.Rotate(math.Pi / 2)
	direction := radToNearestDeg(obj.Direction())
	if direction != 90 {
		t.Errorf("Direction was (%v) should be 90", direction)
	}
}
