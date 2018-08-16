package gologo

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

// TestSimpleObjectCreation : Test basic Object creation by building the object
// then reading values back out to confirm values
// A simple object has no renderer
func TestSimpleObjectCreation(t *testing.T) {
	var defaultX, defaultY, translateX, translateY float32 = 500, 400, 150, 100

	model := mgl32.Vec3{defaultX, defaultY, 0.0}
	obj := CreateObject(model)
	obj.SetZOrder(1)

	// Test default position
	x, y := obj.GetPosition()
	if x != defaultX || y != defaultY {
		t.Errorf("Location was x (%v), y (%v) should be x (%v), y (%v)",
			x, y, defaultX, defaultY)
	}

	// Test translation
	obj.Translate(translateX, translateY)
	x, y = obj.GetPosition()
	if x != defaultX+translateX || y != defaultY+translateY {
		t.Errorf("After translation location was x (%v), y (%v) should be x (%v), y (%v)",
			x, y, defaultX+translateX, defaultY+translateY)
	}

	// Test default direction
	dDirection := obj.Direction()
	if dDirection != 0 {
		t.Errorf("Direction was (%v) should be 0", dDirection)
	}

	// Test rotating by 90 degrees
	obj.Rotate(math.Pi / 2)
	rDirection := radToNearestDeg(obj.Direction())
	if rDirection != 90 {
		t.Errorf("Direction was (%v) should be 90", rDirection)
	}

	// Test object has age
	age := obj.GetAge()
	if age < 0 {
		t.Errorf("Age was (%v) should not be negative", age)
	}
}
