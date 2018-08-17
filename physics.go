package gologo

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

/////////////////////////////////////////////////////////////
// Tick
//

type TickState struct {
	Zero     float64
	Start    float64
	End      float64
	Interval float64
}

func InitTick() error {
	tick.Zero = glfw.GetTime()
	tick.End = tick.Zero

	return nil
}

func GetTime() int {
	return int(1000 * (glfw.GetTime() - tick.Zero))
}

func GetTickTime() int {
	return int(1000 * (tick.End - tick.Zero))
}

func Tick() {
	time := glfw.GetTime()
	tick.Interval = time - tick.End
	tick.End = time
}

/////////////////////////////////////////////////////////////
// Tags
//

func Integrate(duration float64) {
	for _, object := range integrated {
		object.Integrate(duration)
	}
}

func TagIntegrate(object *Object) {
	integrated = append(integrated, object)
}

func UntagIntegrate(object *Object) {
	for i := 0; i < len(integrated); i++ {
		if object == integrated[i] {
			if len(integrated) > 1 {
				integrated = append(integrated[:i], integrated[i+1:]...)
			} else {
				integrated = integrated[0:0]
			}
			i--
		}
	}
}

/////////////////////////////////////////////////////////////
// Physics primitives
//

type Primitive interface {
	InitFromRenderer(r Renderer) error
	IsContainedInRect(obj Object, rect Rect) bool
	OverlapsWithRect(obj Object, rect Rect) bool
	Clone() Primitive
}

type Circle struct {
	Radius float32
}

func (c *Circle) InitFromRenderer(r Renderer) error {
	switch renderer := r.(type) {
	case *MeshRenderer:
		if renderer.VertexCount == 0 {
			return errors.New("object's meshRenderer has no vertices")
		}

		var minX, maxX, minY, maxY float64
		mesh := renderer.MeshVertices

		for i := 0; i < len(mesh); i = i + glMeshStride {
			minX = math.Min(minX, float64(mesh[i]))
			maxX = math.Max(maxX, float64(mesh[i]))
			minY = math.Min(minY, float64(mesh[i+1]))
			maxY = math.Max(maxY, float64(mesh[i+1]))
		}

		size := math.Max(maxX-minX, maxY-minY)
		c.Radius = float32((size * circleMeshSizeFactor) / 2)
	default:
		return errors.New("failed to init renderer from unsupported type: %T")
	}

	return nil
}

func (c *Circle) IsContainedInRect(obj Object, rect Rect) bool {
	x, y := obj.GetPosition()

	xMin, xMax, yMin, yMax := getRectMinMax(rect)

	return y+c.Radius <= yMax &&
		y-c.Radius >= yMin &&
		x+c.Radius <= xMax &&
		x-c.Radius >= xMin
}

func (c *Circle) OverlapsWithRect(obj Object, rect Rect) bool {
	x, y := obj.GetPosition()

	xMin, xMax, yMin, yMax := getRectMinMax(rect)

	return y-c.Radius <= yMax &&
		y+c.Radius >= yMin &&
		x-c.Radius <= xMax &&
		x+c.Radius >= xMin
}

func (c *Circle) Clone() Primitive {
	return &Circle{
		Radius: c.Radius,
	}
}

func CalcCircleCircleContact(
	c1 *Circle,
	c1Pos mgl32.Vec3,
	c2 *Circle,
	c2Pos mgl32.Vec3) (mgl32.Vec3, mgl32.Vec3, float32) {
	v1 := c1Pos.Sub(c2Pos)
	v1Len := v1.Len()
	penetration := (c1.Radius + c2.Radius) - v1Len
	factor := (c2.Radius - penetration/2) / v1Len
	contactPoint := c2Pos.Add(v1.Mul(factor))
	contactNormal := v1.Normalize()

	return contactPoint, contactNormal, penetration
}

/////////////////////////////////////////////////////////////
// Physics Rigid Body
//

// RigidBody : Represents all rigid body physics data for an Object.
// InverseMass is one divided by the mass of the Object.  Representing mass this way allows infinite mass to be represented easily.
// InverseInertia is one divided by the rotational inertia of the Object around the z-axis.  Representing intertia this way allows infinite rotational inertial to be represented easily.
// Velocity is the linear and rotational velocity of the Object.
// LinearDamping is a damping constant applied to any linear acceleration of the Object.  Linear damping can be used to counteract the creation of energy due to simulation inaccuracies.
// AngularDamping is a damping constant applied to any rotational acceleration of the Object.  Angular damping can be used to counteract the creation of energy due to simulation inaccuracies.
// Forces accumulates the linear forces acting on the object during this simulation tick.
// Torques accumulates the rotational forces acting on the object during this simulation tick.
// Accelaration is the linear acceleration this object experienced during the this simulation tick.
type RigidBody struct {
	InverseMass     float64
	InverseInertia  float64
	LinearVelocity  mgl32.Vec3
	AngularVelocity float64
	LinearDamping   float64
	AngularDamping  float64
	Forces          mgl32.Vec3
	Torques         float64
	Acceleration    mgl32.Vec3
}

func (b *RigidBody) Integrate(duration float64) (mgl32.Vec3, float64) {
	// Calculate acceleration due to accumulated forces
	b.Acceleration = b.Forces.Mul(float32(b.InverseMass))
	angularAccel := b.Torques * b.InverseInertia

	// Update linear and angular velocity based on acceleration
	b.LinearVelocity = b.LinearVelocity.Add(b.Acceleration)
	b.AngularVelocity += angularAccel

	// Calculate damping factors
	linearDamping := float32(math.Pow(b.LinearDamping, duration))
	angularDamping := math.Pow(b.AngularDamping, duration)

	// Apply damping based on total linear and angular velocity
	b.LinearVelocity = b.LinearVelocity.Mul(linearDamping)
	b.AngularVelocity = b.AngularVelocity * angularDamping

	return b.LinearVelocity, b.AngularVelocity
}

func (b *RigidBody) AccumulateForceAtLocalPoint(f mgl32.Vec3, p mgl32.Vec3) {
	// Calculate the linear force applied through the centre of gravity
	// This adds the entire force throught the center of gravity.  We might be able to do
	// better.
	b.Forces = b.Forces.Add(f)

	// Calculate the rotational force applied around the z-axis
	// Use the 2D analogue of cross product to calculate the magnitude of the torque
	b.Torques += float64(p.X()*f.Y() - p.Y()*f.X())
}
