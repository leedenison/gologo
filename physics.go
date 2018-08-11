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
// Physics primitives
//

type Primitive interface {
	InitFromRenderer(r Renderer) error
	GetInverseMass() float32
	IsContainedInRect(obj Object, rect Rect) bool
	Clone() Primitive
}

type Circle struct {
	InverseMass float32
	Radius      float32
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
		radius := float32((size * circleMeshSizeFactor) / 2)
		area := math.Pi * radius * radius
		inverseMass := float32(1)

		if area > 0 {
			inverseMass = 1 / (area * areaToMassRatio)
		}
		c.InverseMass = inverseMass
	default:
		return errors.New("failed to init renderer from unsupported type: %T")
	}

	return nil
}

func (c *Circle) GetInverseMass() float32 {
	return c.InverseMass
}

func (c *Circle) IsContainedInRect(obj Object, rect Rect) bool {
	var xMin, xMax, yMin, yMax float32
	x, y := obj.GetPosition()

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

	return y+c.Radius <= yMax &&
		y-c.Radius >= yMin &&
		x+c.Radius <= xMax &&
		x-c.Radius >= xMin
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
