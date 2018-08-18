package gologo

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var createObjTests = []struct {
	posX, posY float32
	zOrder     int
}{
	{200, 300, 1},
	{0, 0, 2},
	{-100, -200, 3},
}

// TestCreateObject : Test basic Object creation by building the object
// then reading values back out to confirm values
// A simple object has no renderer
func TestCreateObject(t *testing.T) {
	var position mgl32.Vec3
	var obj *Object
	var x, y float32
	var age, direction int

	for _, tt := range createObjTests {
		position = mgl32.Vec3{tt.posX, tt.posY, 0.0}
		obj = CreateObject(position)
		obj.SetZOrder(tt.zOrder)

		// Test default position
		x, y = obj.GetPosition()
		if x != tt.posX || y != tt.posY {
			t.Errorf("Location was x (%v), y (%v) should be x (%v), y (%v)",
				x, y, tt.posX, tt.posY)
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
	}
}

/*
// TestObject_Translate : Test that Object translation moves the object by the specified amount
func TestObject_Translate(t *testing.T) {
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

// TestObject_Rotate : Test that Object rotation changes the orientation of the object by the specified amount
func TestObject_Rotate(t *testing.T) {
	position := mgl32.Vec3{positionX, positionY, 0.0}
	obj := CreateObject(position)

	// Test rotating by 90 degrees
	obj.Rotate(math.Pi / 2)
	direction := radToNearestDeg(obj.Direction())
	if direction != 90 {
		t.Errorf("Direction was (%v) should be 90", direction)
	}
}
*/
var CasesObject_Integrate = []struct {
	name     string
	duration float64
	in       Object
	out      Object
}{
	{
		name:     "no velocity no forces no damping",
		duration: 1.0,
		in: Object{
			Position:    mgl32.Vec3{0, 0, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     1.0,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{0.0, 0.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   1.0,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{0.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
		out: Object{
			Position:    mgl32.Vec3{0, 0, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     1.0,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{0.0, 0.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   1.0,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{0.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
	},
	{
		name:     "constant linear velocity no forces no damping",
		duration: 1.0,
		in: Object{
			Position:    mgl32.Vec3{0, 0, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     1.0,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{10.0, 10.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   1.0,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{0.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
		out: Object{
			Position:    mgl32.Vec3{10, 10, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     1.0,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{10.0, 10.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   1.0,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{0.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
	},
	{
		name:     "deceleration force no damping",
		duration: 1.0,
		in: Object{
			Position:    mgl32.Vec3{0, 0, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     0.5,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{10.0, 10.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   1.0,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{-10.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
		out: Object{
			Position:    mgl32.Vec3{5, 10, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     0.5,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{5.0, 10.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   1.0,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{-10.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
	},
	{
		name:     "deceleration force 5 percent damping",
		duration: 1.0,
		in: Object{
			Position:    mgl32.Vec3{0, 0, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     0.5,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{10.0, 10.0, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   0.95,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{-10.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
		out: Object{
			Position:    mgl32.Vec3{4.75, 9.5, 0},
			Orientation: 0,
			Body: &RigidBody{
				InverseMass:     0.5,
				InverseInertia:  1.0,
				LinearVelocity:  mgl32.Vec3{4.75, 9.5, 0.0},
				AngularVelocity: 0.0,
				LinearDamping:   0.95,
				AngularDamping:  1.0,
				Forces:          mgl32.Vec3{-10.0, 0.0, 0.0},
				Torques:         0.0,
			},
		},
	},
}

// TestObject_Integrate : Test that integration correctly updates position and orientation
func TestObject_Integrate(t *testing.T) {
	for _, tc := range CasesObject_Integrate {
		t.Run(tc.name, func(t *testing.T) {
			tc.in.Integrate(tc.duration)
			if tc.in.Position != tc.out.Position {
				t.Errorf("Expected post-integration position: %v, found %v",
					tc.out.Position, tc.in.Position)
			}
			if tc.in.Orientation != tc.out.Orientation {
				t.Errorf("Expected post-integration orientation: %v, found %v",
					tc.out.Orientation, tc.in.Orientation)
			}
			if tc.in.Body.LinearVelocity != tc.out.Body.LinearVelocity {
				t.Errorf("Expected post-integration linear velocity: %v, found %v",
					tc.out.Body.LinearVelocity, tc.in.Body.LinearVelocity)
			}
			if tc.in.Body.AngularVelocity != tc.out.Body.AngularVelocity {
				t.Errorf("Expected post-integration angular velocity: %v, found %v",
					tc.out.Body.AngularVelocity, tc.in.Body.AngularVelocity)
			}
			if tc.in.Body.InverseMass != tc.out.Body.InverseMass ||
				tc.in.Body.InverseInertia != tc.out.Body.InverseInertia ||
				tc.in.Body.LinearDamping != tc.out.Body.LinearDamping ||
				tc.in.Body.AngularDamping != tc.out.Body.AngularDamping {
				t.Errorf("Integration modified rigid body constant - expected: %v, found %v",
					tc.out.Body, tc.in.Body)
			}
		})
	}
}
