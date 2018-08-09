package gologo

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
	InitFromMesh(mesh []float32) Primitive
	GetInverseMass() float32
	IsOnScreen(x float32, y float32) bool
	Clone() Primitive
}

type Circle struct {
	InverseMass float32
	Radius      float32
}

func (c *Circle) InitFromMesh(mesh []float32) Primitive {
	var minX, maxX, minY, maxY float64

	for i := 0; i < len(mesh); i = i + GL_MESH_STRIDE {
		minX = math.Min(minX, float64(mesh[i]))
		maxX = math.Max(maxX, float64(mesh[i]))
		minY = math.Min(minY, float64(mesh[i+1]))
		maxY = math.Max(maxY, float64(mesh[i+1]))
	}

	size := math.Max(maxX-minX, maxY-minY)
	radius := float32((size * CIRCLE_MESH_SIZE_FACTOR) / 2)
	area := math.Pi * radius * radius
	inverseMass := float32(1)

	if area > 0 {
		inverseMass = 1 / (area * AREA_TO_MASS_RATIO)
	}

	return &Circle{
		InverseMass: inverseMass,
		Radius:      radius,
	}
}

func (c *Circle) GetInverseMass() float32 {
	return c.InverseMass
}

func (c *Circle) IsOnScreen(x float32, y float32) bool {
	return y+c.Radius <= float32(windowState.Height) &&
		y-c.Radius >= 0.0 &&
		x+c.Radius <= float32(windowState.Width) &&
		x-c.Radius >= 0.0
}

func (c *Circle) Clone() Primitive {
	return &Circle{
		InverseMass: c.InverseMass,
		Radius:      c.Radius,
	}
}

func CalcCircleCircleContact(
	c1 *Circle,
	c1Model mgl32.Mat4,
	c2 *Circle,
	c2Model mgl32.Mat4) (mgl32.Vec4, mgl32.Vec4, float32) {
	v1 := c1Model.Col(3).Vec3().Sub(c2Model.Col(3).Vec3())
	v1Len := v1.Len()
	penetration := (c1.Radius + c2.Radius) - v1Len
	factor := (c2.Radius - penetration/2) / v1Len
	contactPoint := c2Model.Col(3).Add(v1.Mul(factor).Vec4(1.0))
	contactNormal := v1.Normalize().Vec4(1.0)

	return contactPoint, contactNormal, penetration
}

func IsOnScreen(x float32, y float32) bool {
	return y <= float32(windowState.Height) &&
		y >= 0.0 &&
		x <= float32(windowState.Width) &&
		x >= 0.0
}
