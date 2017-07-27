package gologo

import (
    "github.com/go-gl/mathgl/mgl32"
    "math"
)

type Primitive interface {
    GetInverseMass() float32
}

type Circle struct {
    InverseMass float32
    Radius float32
}

func (c *Circle) GetInverseMass() float32 {
    return c.InverseMass
}

func CalcCircleCircleContact(
        c1 *Circle,
        c1Model mgl32.Mat4,
        c2 *Circle,
        c2Model mgl32.Mat4) (mgl32.Vec4, mgl32.Vec4, float32) {
    v1 := c1Model.Col(3).Vec3().Sub(c2Model.Col(3).Vec3())
    v1Len := v1.Len()
    penetration := (c1.Radius + c2.Radius) - v1Len
    factor := (c2.Radius - penetration / 2) / v1Len
    contactPoint := c2Model.Col(3).Add(v1.Mul(factor).Vec4(1.0))
    contactNormal := v1.Normalize().Vec4(1.0)

    return contactPoint, contactNormal, penetration
}

func InitCircleFromMesh(mesh []float32) Primitive {
    var minX, maxX, minY, maxY float64

    for i := 0; i < len(mesh); i = i + GL_MESH_STRIDE {
        minX = math.Min(minX, float64(mesh[i]))
        maxX = math.Max(maxX, float64(mesh[i]))
        minY = math.Min(minY, float64(mesh[i + 1]))
        maxY = math.Max(maxY, float64(mesh[i + 1]))
    }

    size := math.Max(maxX - minX, maxY - minY)
    radius := float32((size * CIRCLE_MESH_SIZE_FACTOR) / 2)
    area := math.Pi * radius * radius
    inverseMass := float32(1)

    if area > 0 {
        inverseMass = 1 / (area * AREA_TO_MASS_RATIO)
    }

    return &Circle {
        InverseMass: inverseMass,
        Radius: radius,
    }
}

func CircleIsOnScreen(circle *Circle, model mgl32.Mat4) bool {
    position := model.Col(3)

    return position.Y() - circle.Radius <= float32(glWin.Height) &&
        position.Y() + circle.Radius >= 0.0 &&
        position.X() - circle.Radius <= float32(glWin.Width) &&
        position.X() + circle.Radius >= 0.0
}
