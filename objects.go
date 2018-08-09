package gologo

import (
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// Object : Struct to hold fundamental object for gologo
// Model is the mgl32 model
// ZOrder is the gologo managed height order of the objects - 0 is valid
// Creation is a automatically managed time the object was created
// Primitive is the physics primitive for this object - can be nil
// Renderer is the gl renderer for this object - can be nil
type Object struct {
	Model     mgl32.Mat4
	ZOrder    int
	Creation  int
	Primitive Primitive
	Renderer  Renderer
}

// CreateObject Main function to create a standard object
// required a model to create an object
func CreateObject(model mgl32.Mat4) *Object {
	return &Object{
		Model:    model,
		Creation: GetTickTime(),
	}
}

// GetModel : Returns the model for this object
func (o *Object) GetModel() mgl32.Mat4 {
	return o.Model
}

// GetAge : Returns age of object since creation
func (o *Object) GetAge() int {
	return GetTickTime() - o.Creation
}

// GetPosition : Returns X and Y co-ords of object centre in 2D
func (o *Object) GetPosition() (float32, float32) {
	return o.Model.Col(3).X(), o.Model.Col(3).Y()
}

// SetPosition : Sets X and Y co-ords of object centre in 2D
func (o *Object) SetPosition(x float32, y float32) {
	o.Model.SetCol(3, mgl32.Vec4{x, y, 0.0, 1.0})
}

// SetPositionVec2 : Sets X and Y co-ords of object centre
// in 2D using a vector of 2 elements
func (o *Object) SetPositionVec2(p mgl32.Vec2) {
	o.Model.SetCol(3, p.Vec4(0.0, 1.0))
}

// Translate : Move the position by the supplied X and Y values
// relative to the current position
func (o *Object) Translate(x float32, y float32) {
	o.Model = mgl32.Translate3D(x, y, 0.0).Mul4(o.Model)
}

// Direction : Return the angle the object has been rotated since it was created
func (o *Object) Direction() float64 {
	return math.Atan2(float64(o.Model.At(1, 1)), float64(o.Model.At(0, 1))) - math.Pi/2
}

// DirectionNormal : Return the normal to the angle the object has been rotated since it was created
func (o *Object) DirectionNormal() mgl32.Vec3 {
	return o.Model.Col(1).Vec3().Normalize()
}

// DirectionOf : Calculates the direction in radians to the passed in object
// from the receiving object
func (o *Object) DirectionOf(other *Object) float64 {
	direction := other.Model.Col(3).Vec3().Sub(o.Model.Col(3).Vec3())

	return math.Atan2(float64(direction[1]), float64(direction[0])) - math.Pi/2
}

// Rotate : Rotate the object by the supplied angle in radians
func (o *Object) Rotate(angle float32) {
	rotation := mgl32.HomogRotate3DZ(angle)

	o.Model = o.Model.Mul4(rotation)
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
	if primitive.GetInverseMass() == 0 {
		// Primitive hasn't been set up, initialise it with defaults and
		// store it in the object and ignore the clone var as it's new
		return o.initialisePrimitive(primitive)
	}

	if clone {
		o.Primitive = primitive.Clone()
	} else {
		o.Primitive = primitive
	}

	return nil
}

// initialisePrimitive : Creates a default primitive for this object.
// Will use the mesh to calculate the primitive
// Will return an error if the Renderer is not set, the Renderer is not
// of type MeshRenderer, or has no vertices
func (o *Object) initialisePrimitive(primitive Primitive) error {
	if o.Renderer == nil {
		return errors.New("object has no renderer")
	}

	meshRenderer, ok := o.Renderer.(*MeshRenderer)
	if !ok {
		return errors.New("object's Renderer wouldn't cast to MeshRenderer, is it a MeshRenderer?")
	}

	if meshRenderer.VertexCount == 0 {
		return errors.New("object renderer has no vertices")
	}

	o.Primitive = primitive.InitFromMesh(meshRenderer.MeshVertices)

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
