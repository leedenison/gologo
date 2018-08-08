package gologo

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

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
