package gologo

import "math"

type Vector2 struct {
    x, y   float64
}

func (v *Vector2) Scale(a float64) {
    v.x = v.x * a
    v.y = v.y * a
}

func (v *Vector2) Add(a *Vector2) {
    v.x += a.x
    v.y += a.y
}

func (v *Vector2) AddAngle(a *Vector2) {
    // cos(a + b) = cos(a).cos(b) - sin(a).sin(b)
    v.x = v.x * a.x - v.y * a.y

    // sin(a + b) = sin(a).cos(b) + cos(a).sin(b)
    v.y = v.y * a.x + v.x * a.y
}

func (v *Vector2) Subtract(a *Vector2) {
    v.x -= a.x
    v.y -= a.y
}

func (v *Vector2) AddScaledVector(a *Vector2, b float64) {
    v.x += a.x * b
    v.y += a.y * b
}

func UnitVector(v *Vector2) *Vector2 {
    magnitude := Magnitude(v)
    return &Vector2 { v.x / magnitude, v.y / magnitude }
}

func SquareMagnitude(v *Vector2) float64 {
    return math.Pow(v.x, 2) + math.Pow(v.y, 2)
}

func Magnitude(v *Vector2) float64 {
    return math.Sqrt(SquareMagnitude(v))
}

func ComponentProduct(a *Vector2, b *Vector2) *Vector2 {
    return &Vector2 { a.x * b.x, a.y * b.y }
}

func DotProduct(a *Vector2, b *Vector2) float64 {
    return a.x * b.x + a.y * b.y
}

func Orthogonal(v *Vector2) *Vector2 {
    return &Vector2 { v.y, -v.x }
}

func DirectionVector(from *Vector2, to *Vector2) *Vector2 {
    return &Vector2 { to.x - from.x, to.y - from.y }
}

func UnitDirectionVector(from *Vector2, to *Vector2) *Vector2 {
    return UnitVector(DirectionVector(from, to))
}

