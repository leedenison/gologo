package gologo

import (
	"image/color"
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

func Colors(gradients []float64, count int) []color.RGBA {
	var stride = 4
	var red = 1
	var green = 2
	var blue = 3

	result := []color.RGBA{}

	for i := 0; i < count; i++ {
		q := float64(i) / float64(count)

		for j := stride; j < len(gradients); j += stride {
			if q < gradients[j] {
				p := q - gradients[j-stride]
				r := uint8((gradients[j+red]-gradients[j-stride+red])*p + gradients[j-stride+red])
				g := uint8((gradients[j+green]-gradients[j-stride+green])*p + gradients[j-stride+green])
				b := uint8((gradients[j+blue]-gradients[j-stride+blue])*p + gradients[j-stride+blue])
				result = append(result, color.RGBA{r, g, b, 255})
				break
			}
		}
	}

	return result
}
