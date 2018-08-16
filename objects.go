package gologo

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

// Object : Struct to hold fundamental object for gologo
// Position is the position of the object origin in world space
// Orientation is the rotation of the object in radians
// ZOrder is the gologo managed height order of the objects - 0 is valid
// Creation is a automatically managed time the object was created
// Primitive is the physics primitive for this object - can be nil
// Renderer is the gl renderer for this object - can be nil
// Body is the RigidBody physics data for this object - can be nil
type Object struct {
	Position    mgl32.Vec3
	Orientation float64
	Scale       float64
	ZOrder      int
	Creation    int
	Primitive   Primitive
	Renderer    Renderer
	Body        *RigidBody
}

// CreateObject Main function to create a standard object
// requires a position to create an object
func CreateObject(position mgl32.Vec3) *Object {
	return &Object{
		Position: position,
		Creation: GetTickTime(),
	}
}

// GetModel : Returns the model for this object
func (o *Object) GetModel() mgl32.Mat4 {
	translate := mgl32.Translate3D(o.Position.X(), o.Position.Y(), o.Position.Z())
	scale := mgl32.Scale3D(float32(o.Scale), float32(o.Scale), 1.0)
	rotate := mgl32.HomogRotate3DZ(float32(o.Orientation))
	return translate.Mul4(rotate.Mul4(scale))
}

// WorldSpace : Returns the world space point corresponding to the supplied object space co-ordinate
func (o *Object) WorldSpace(c mgl32.Vec3) mgl32.Vec3 {
	return o.GetModel().Mul4x1(c.Vec4(1.0)).Vec3()
}

// Integrate : Updates RigidBody velocity and then updates the Model matrix with the result
func (o *Object) Integrate(duration float64) {
	linear, angular := o.Body.Integrate(duration)
	o.Position = o.Position.Add(linear)
	// Normalize the orientation to be within the first 2 pi radians
	o.Orientation = math.Mod(o.Orientation+angular, math.Pi*2)
}

// GetAge : Returns age of object since creation
func (o *Object) GetAge() int {
	return GetTickTime() - o.Creation
}

// GetPosition : Returns X and Y co-ords of object centre in 2D
func (o *Object) GetPosition() (float32, float32) {
	return o.Position.X(), o.Position.Y()
}

// SetPosition : Sets X and Y co-ords of object centre in 2D
func (o *Object) SetPosition(x float32, y float32) {
	o.Position = mgl32.Vec3{x, y, 0.0}
}

// SetPositionVec2 : Sets X and Y co-ords of object centre
// in 2D using a vector of 2 elements
func (o *Object) SetPositionVec2(p mgl32.Vec2) {
	o.Position = p.Vec3(0.0)
}

// Translate : Move the position by the supplied X and Y values
// relative to the current position
func (o *Object) Translate(x float32, y float32) {
	o.Position = o.Position.Add(mgl32.Vec3{x, y, 0.0})
}

// Direction : Return the angle the object has been rotated since it was created
func (o *Object) Direction() float64 {
	return o.Orientation
}

// DirectionVector : Return the normalized vector in the direction the object has been rotated
func (o *Object) DirectionVector() mgl32.Vec3 {
	x := float32(math.Cos(o.Orientation))
	y := float32(math.Sin(o.Orientation))
	return mgl32.Vec3{x, y, 0.0}
}

// DirectionNormal : Return the normal to the angle the object has been rotated since it was created
func (o *Object) DirectionNormal() mgl32.Vec3 {
	// Normal to a vector [x, y] in 2D is [-y, x]
	directionVector := o.DirectionVector()
	return mgl32.Vec3{-directionVector.Y(), directionVector.X(), 0.0}
}

// DirectionOf : Calculates the direction in radians to the passed in object
// from the receiving object
func (o *Object) DirectionOf(other *Object) float64 {
	direction := other.Position.Sub(o.Position).Normalize()

	return math.Acos(float64(direction[0]))
}

// Rotate : Rotate the object by the supplied angle in radians
func (o *Object) Rotate(angle float32) {
	o.Orientation = math.Mod(o.Orientation+float64(angle), math.Pi*2)
}

// SetZOrder : Sets the height of the object in 3D space
// as an integer compared with other objects
func (o *Object) SetZOrder(z int) {
	o.ZOrder = z
}

// GetRenderer : Returns the renderer for this object
func (o *Object) GetRenderer() Renderer {
	return o.Renderer
}

// SetRenderer : Sets the renderer for this object, optionally cloning it
func (o *Object) SetRenderer(renderer Renderer, clone bool) {
	if clone {
		o.Renderer = renderer.Clone()
	} else {
		o.Renderer = renderer
	}
}

// GetPrimitive : Returns the primitive for this object
func (o *Object) GetPrimitive() Primitive {
	return o.Primitive
}

// SetPrimitive : Sets the primitive for this object, optionally cloning it
// and deriving default values if needed
func (o *Object) SetPrimitive(primitive Primitive, clone bool) error {
	if clone {
		o.Primitive = primitive.Clone()
	} else {
		o.Primitive = primitive
	}

	return nil
}

// InitialisePrimitive : Creates a default primitive for this object.
// Will use the mesh to calculate the primitive
// Will return an error if the Renderer is not set, the Renderer is not
// of type MeshRenderer, or has no vertices
func (o *Object) InitialisePrimitive() error {
	if o.Renderer == nil {
		return errors.New("object has no renderer")
	}

	err := o.Primitive.InitFromRenderer(o.Renderer)

	if err != nil {
		return errors.Wrap(err, "failed to initialise primitive from renderer")
	}

	return nil
}

// Clone : creates a distinct copy of the receiving object
func (o *Object) Clone() *Object {
	objectCopy := *o

	if o.Renderer != nil {
		objectCopy.Renderer = o.Renderer.Clone()
	}

	if o.Primitive != nil {
		objectCopy.Primitive = o.Primitive.Clone()
	}

	return &objectCopy
}

// IsContainedInRect : returns true if the object
// A primitive is in the rect if the entire primitive is
// A Mesh without a primitive is in the rect if it's centre is
func (o *Object) IsContainedInRect(rect Rect) bool {
	var xMin, xMax, yMin, yMax float32

	if o.GetPrimitive() != nil {
		return o.GetPrimitive().IsContainedInRect(*o, rect)
	}

	x, y := o.GetPosition()

	if rect[0][0] > rect[1][0] {
		xMin = rect[1][0]
		xMax = rect[0][0]
	} else {
		xMin = rect[0][0]
		xMax = rect[1][0]
	}

	if rect[0][1] > rect[1][1] {
		yMin = rect[1][1]
		yMax = rect[0][1]
	} else {
		yMin = rect[0][1]
		yMax = rect[1][1]
	}

	return y <= yMax &&
		y >= yMin &&
		x <= xMax &&
		x >= xMin
}

// ByZOrder : Height ordering array for Objects
type ByZOrder []*Object

func (s ByZOrder) Len() int {
	return len(s)
}

func (s ByZOrder) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByZOrder) Less(i int, j int) bool {
	return s[i].ZOrder < s[j].ZOrder
}
