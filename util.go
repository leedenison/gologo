package gologo

import "math"

func Clamp(value, limit1, limit2 float64) float64 {
    var min, max float64
    if limit1 > limit2 {
        min = limit2
        max = limit1
    } else {
        min = limit1
        max = limit2
    }

    switch {
    case value < min:
        return min
    case value > max:
        return max
    }
    return value
}

func LineGradientIntercept(a *Vector2, b *Vector2) (float64, float64) {
    deltaX := a.x - b.x
    deltaY := a.y - b.y

    if deltaX == 0 && deltaY == 0 {
        // The two points are the same
        return math.NaN(), math.NaN()
    } else if deltaX == 0 {
        // The line is vertical
        return math.Inf(1), math.NaN()
    } else if deltaY == 0 {
        // The line is horizontal
        return 0, float64(a.y)
    } else {
        gradient := float64(deltaY) / float64(deltaX)
        return gradient, float64(a.y) - gradient * float64(a.x)
    }
}

func ReflectVector(unitNormal *Vector2, v *Vector2) Vector2 {
    dot := DotProduct(v, unitNormal)

    // Reflect the velocity around the normal vector
    return Vector2 {
        x: v.x - 2 * dot * unitNormal.x,
        y: v.y - 2 * dot * unitNormal.y,
    }
}

func CircleArea(radius float64) float64 {
    return math.Pi * math.Pow(radius, 2)
}