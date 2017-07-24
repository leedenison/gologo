package gologo

import (
    "github.com/go-gl/mathgl/mgl32"
)

type Path struct {
    Segments [][4]mgl32.Vec2
    Times []float32
}

func (p *Path) GetPosition(time float32) (int, int) {
    var accumulator float32

    for i := 0; i < len(p.Times); i++ {
        if accumulator + p.Times[i] >= time {
            segmentTime := (time - accumulator) / p.Times[i]
            position := mgl32.CubicBezierCurve2D(
                segmentTime,
                p.Segments[i][0],
                p.Segments[i][1],
                p.Segments[i][2],
                p.Segments[i][3])

            return int(position.X()), int(position.Y())
        }

        accumulator += p.Times[i]
    }

    return 0, 0
}

func CreatePath(path []int) *Path {
    if len(path) < 4 || len(path) % 2 != 0 {
        panic("Path must an even number of integers")
    }

    points := convertToVec2(path)
    result := Path {}

    if len(points) - 1 == 1 {
        segment := calcSingleSegment(points)
        result.Segments = append(result.Segments, segment)
        result.Times = []float32 { 1.0 }
    } else {
        coefficients, input := calcInput(points)
        controls := solveTriDiagonal(coefficients, input)
        result.Segments = calcSegments(points, controls)

        result.Times = make([]float32, len(result.Segments))

        calcSegmentLengths(result)
    }

    return &result
}

func calcSegmentLengths(path Path) {
    var totalLength float32

    for i := 0; i < len(path.Segments); i++ {
        s := path.Segments[i]
        chord := s[3].Sub(s[0]).Len()
        controlNet := s[1].Sub(s[0]).Len() + s[2].Sub(s[1]).Len() + s[3].Sub(s[2]).Len()
        path.Times[i] = (chord + controlNet) / 2
        totalLength += path.Times[i]
    }

    for i := 0; i < len(path.Segments); i++ {
        path.Times[i] = path.Times[i] / totalLength
    }
}

func calcSingleSegment(points []mgl32.Vec2) [4]mgl32.Vec2 {
    control1 := mgl32.Vec2 {
        (2 * points[0].X() + points[1].X()) / 3,
        (2 * points[0].Y() + points[1].Y()) / 3,
    }

    control2 := mgl32.Vec2 {
        2 * control1.X() - points[0].X(),
        2 * control1.Y() - points[0].Y(),
    }

    segment := [4]mgl32.Vec2 {
        mgl32.Vec2 { points[0].X(), points[0].X() },
        control1,
        control2,
        mgl32.Vec2 { points[1].Y(), points[1].Y() },
    }

    return segment
}

func convertToVec2(points []int) []mgl32.Vec2 {
    result := []mgl32.Vec2 {}

    for i := 0; i < len(points); i = i + 2 {
        result = append(result, mgl32.Vec2 {
            float32(points[i]),
            float32(points[i + 1]),
        })
    }

    return result
}

func calcSegments(points []mgl32.Vec2, bControls []mgl32.Vec2) [][4]mgl32.Vec2 {
    result := [][4]mgl32.Vec2 {}

    for i := 0; i < len(points) - 1; i++ {
        var cubic1, cubic2 mgl32.Vec2

        if i == 0 {
            cubic1 = mgl32.Vec2 {
                (2.0 * points[i].X() + bControls[i].X()) / 3.0,
                (2.0 * points[i].Y() + bControls[i].Y()) / 3.0,
            }
            cubic2 = mgl32.Vec2 {
                2.0 * cubic1.X() - points[i].X(),
                2.0 * cubic1.Y() - points[i].Y(),
            }
        } else if i == len(points) - 2 {
            cubic1 = mgl32.Vec2 {
                (2.0 * bControls[i - 1].X() + points[i + 1].X()) / 3.0,
                (2.0 * bControls[i - 1].Y() + points[i + 1].Y()) / 3.0,
            }
            cubic2 = mgl32.Vec2 {
                2.0 * cubic1.X() - bControls[i - 1].X(),
                2.0 * cubic1.Y() - bControls[i - 1].Y(),
            }
        } else {
            cubic1 = mgl32.Vec2 {
                (2.0 * bControls[i - 1].X() + bControls[i].X()) / 3.0,
                (2.0 * bControls[i - 1].Y() + bControls[i].Y()) / 3.0,
            }
            cubic2 = mgl32.Vec2 {
                2.0 * cubic1.X() - bControls[i - 1].X(),
                2.0 * cubic1.Y() - bControls[i - 1].Y(),
            }
        }

        result = append(result, [4]mgl32.Vec2 { points[i], cubic1, cubic2, points[i + 1] })
    }

    return result
}

func solveTriDiagonal(coefficients [][3]float32, input []mgl32.Vec2) []mgl32.Vec2 {
    output := make([]mgl32.Vec2, len(input))

    coefficients[0][2] = coefficients[0][2] / coefficients[0][1]
    output[0] = mgl32.Vec2 {
        input[0].X() / coefficients[0][1],
        input[0].Y() / coefficients[0][1],
    }

    for i := 1; i < len(input); i++ {
        m := 1.0 / (coefficients[i][1] - coefficients[i][0] * coefficients[i - 1][2])
        coefficients[i][2] = m * coefficients[i][2]

        output[i] = mgl32.Vec2 {
            (input[i].X() - coefficients[i][0] * output[i - 1].X()) * m,
            (input[i].Y() - coefficients[i][0] * output[i - 1].Y()) * m,
        }
    }

    for i := len(input) - 2; i >= 0; i-- {
        output[i] = mgl32.Vec2 {
            output[i].X() - coefficients[i][2] * output[i + 1].X(),
            output[i].Y() - coefficients[i][2] * output[i + 1].Y(),
        }
    }

    return output
}

func calcInput(points []mgl32.Vec2) ([][3]float32, []mgl32.Vec2) {
    coefficients := [][3]float32 {}
    input := []mgl32.Vec2 {}

    for i := 1; i < len(points) - 1; i++ {
        if i == 1 {
            coefficients = append(coefficients, [3]float32 { 0, 4, 1 } )
            input = append(input, mgl32.Vec2 {
                6.0 * points[i].X() - points[i - 1].X(),
                6.0 * points[i].Y() - points[i - 1].Y(),
            })
        } else if i == len(points) - 2 {
            coefficients = append(coefficients, [3]float32 { 1, 4, 0 } )
            input = append(input, mgl32.Vec2 {
                6.0 * points[i].X() - points[i + 1].X(),
                6.0 * points[i].Y() - points[i + 1].Y(),
            })
        } else {
            coefficients = append(coefficients, [3]float32 { 1, 4, 1 } )
            input = append(input, mgl32.Vec2 {
                6.0 * points[i].X(),
                6.0 * points[i].Y(),
            })
        }
    }

    return coefficients, input
}
