package gologo

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Rect [2][2]float32

func containsInt(s []int, v int) bool {
	for _, c := range s {
		if c == v {
			return true
		}
	}
	return false
}

func radToNearestDeg(angle float64) int {
	if angle <= -0 {
		angle += 2 * math.Pi
	}
	return int(mgl32.RadToDeg(float32(angle)))
}

func getRectMinMax(rect Rect) (float32, float32, float32, float32) {
	var xMin, xMax, yMin, yMax float32

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

	return xMin, xMax, yMin, yMax
}
