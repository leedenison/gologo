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

func LineGradientIntercept(a Vector, b Vector) (float64, float64) {
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

func UnitDirectionVector(from Vector, to Vector) Vector {
    dirVector := Vector { to.x - from.x, to.y - from.y }
    magnitude := math.Sqrt(math.Pow(dirVector.x, 2) + math.Pow(dirVector.y, 2))
    dirVector.x = dirVector.x / magnitude
    dirVector.y = dirVector.y / magnitude
    return dirVector
}

func DotProduct(a Vector, b Vector) float64 {
    return a.x * b.x + a.y * b.y
}
