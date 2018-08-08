package gologo

import (
	"errors"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

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

// SetZOrder : Sets the height of the object in 3D space
// as an integer compared with other objects
func (o *Object) SetZOrder(z int) {
	o.ZOrder = z
}

func (o *Object) Direction() float64 {
	return math.Atan2(float64(o.Model.At(1, 1)), float64(o.Model.At(0, 1))) - math.Pi/2
}

// DirectionOf : Calculates the direction in radians to the passed in object
// from the receiving object
func (o *Object) DirectionOf(other *Object) float64 {
	direction := other.Model.Col(3).Vec3().Sub(o.Model.Col(3).Vec3())

	return math.Atan2(float64(direction[1]), float64(direction[0])) - math.Pi/2
}

func (o *Object) Rotate(angle float32) {
	rotation := mgl32.HomogRotate3DZ(angle)

	o.Model = o.Model.Mul4(rotation)
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
func (o *Object) SetPrimitive(primitive Primitive, clone bool) {
	if clone {
		o.Primitive = primitive.Clone()
	} else {
		o.Primitive = primitive
	}
}

// SetDefaultPrimitive : Creates a default primitive for this object
// The default is Circle currently and will calculate the circle size
// from the renderers mesh. Will return an error if the Renderer is
// not set, the Renderer is not of type MeshRenderer, or has no vertices
func (o *Object) SetDefaultPrimitive() error {
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

	o.Primitive = InitCircleFromMesh(meshRenderer.MeshVertices)

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
