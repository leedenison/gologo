package gologo

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo/fakes"
	"github.com/leedenison/gologo/opengl"
)

// Need to decide on corner cases like all 0s
var circleContainedInTests = []struct {
	name                                       string
	botX, topX, botY, topY, radius, posX, posY float32
	expResult                                  bool
}{
	{"contained in middle", 0, 400, 0, 400, 50, 200, 200, true},
	{"inside and overlap top right", 0, 400, 0, 400, 50, 380, 380, false},
	{"outside and overlap top right", 0, 400, 0, 400, 50, 425, 425, false},
	{"outside and non-overlap top right", 0, 400, 0, 400, 50, 475, 475, false},
	{"outside and non-overlap top right small", 0, 400, 0, 400, 20, 425, 425, false},
	{"touching right", 0, 400, 0, 400, 100, 500, 400, false},
	{"one unit non-overlap right", 0, 400, 0, 400, 100, 501, 400, false},
}

// TestCircleContainedIn : Test circle contained in rect
func TestCircleContainedIn(t *testing.T) {
	var x, y float32
	var rect Rect
	var position mgl32.Vec3
	var circle Circle
	var isContained bool

	obj := CreateObject(position)

	for _, tc := range circleContainedInTests {
		t.Run(tc.name, func(t *testing.T) {
			obj.SetPosition(tc.posX, tc.posY)
			circle.Radius = tc.radius
			x, y = obj.GetPosition()
			rect = Rect{{tc.botX, tc.botY}, {tc.topX, tc.topY}}

			isContained = circle.IsContainedInRect(*obj, rect)
			if isContained != tc.expResult {
				t.Errorf("IsContInRect %v (exp: %v) with x (%v), y (%v), rad (%v) in x (%v-%v) and y (%v-%v)",
					isContained, tc.expResult, x, y, circle.Radius, tc.botX, tc.topX, tc.botY, tc.topY)
			}
		})
	}
}

var circleOverlapTests = []struct {
	name                                       string
	botX, topX, botY, topY, radius, posX, posY float32
	expResult                                  bool
}{
	{"contained in middle", 0, 400, 0, 400, 50, 200, 200, true},
	{"inside and overlap top right", 0, 400, 0, 400, 50, 380, 380, true},
	{"outside and overlap top right", 0, 400, 0, 400, 50, 425, 425, true},
	{"outside and non-overlap top right", 0, 400, 0, 400, 50, 475, 475, false},
	{"outside and non-overlap top right small", 0, 400, 0, 400, 20, 425, 425, false},
	{"touching right", 0, 400, 0, 400, 100, 500, 400, true},
	{"one unit non-overlap right", 0, 400, 0, 400, 100, 501, 400, false},
}

// TestCircleOverlap : Test circle overlaps with rect
func TestCircleOverlap(t *testing.T) {
	var x, y float32
	var rect Rect
	var position mgl32.Vec3
	var circle Circle
	var overlaps bool

	obj := CreateObject(position)

	for _, tc := range circleOverlapTests {
		t.Run(tc.name, func(t *testing.T) {
			obj.SetPosition(tc.posX, tc.posY)
			circle.Radius = tc.radius
			x, y = obj.GetPosition()
			rect = Rect{{tc.botX, tc.botY}, {tc.topX, tc.topY}}

			overlaps = circle.OverlapsWithRect(*obj, rect)
			if overlaps != tc.expResult {
				t.Errorf("overlapWithRect %v (exp: %v) with x (%v), y (%v), rad (%v) in x (%v-%v) and y (%v-%v)",
					overlaps, tc.expResult, x, y, circle.Radius, tc.botX, tc.topX, tc.botY, tc.topY)
			}
		})
	}
}

var circleFromRendererTests = []struct {
	name      string
	vertices  []float32
	expRadius float32
}{
	{"200 square at origin", []float32{-100, -100, 0, 0, 1, 100, 100, 0, 1, 0,
		-100, 100, 0, 0, 0, -100, -100, 0, 0, 1,
		100, -100, 0, 1, 1, 100, 100, 0, 1, 0}, 65,
	},
	{"50 triangle top right of origin", []float32{0, 0, 0, 0, 1, 50, 50, 0, 1, 0,
		0, 50, 0, 0, 0, 0, 0, 0, 0, 1}, 16.25,
	},
	{"1 square", []float32{0, 0, 0, 0, 1, 1, 1, 0, 1, 0,
		0, 1, 0, 0, 0, 0, 0, 0, 0, 1,
		1, 0, 0, 1, 1, 1, 1, 0, 1, 0}, 0.325,
	},
	{"0 square", []float32{0, 0, 0, 0, 1, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 1,
		0, 0, 0, 1, 1, 0, 0, 0, 1, 0}, 0,
	},
}

// TestCircleFromRenderer : generates circle primitives from
// a mocked renderer which only contains vertices
func TestCircleFromRenderer(t *testing.T) {
	var circle Circle
	var renderer *opengl.MeshRenderer
	var err error

	opengl.CreateMeshRenderer = fakes.CreateVerticesOnlyMeshRendererImpl

	for _, tc := range circleFromRendererTests {
		t.Run(tc.name, func(t *testing.T) {
			renderer, err = opengl.CreateMeshRenderer("", "", []int{}, map[int]interface{}{}, tc.vertices)
			circle.InitFromRenderer(renderer)

			if circle.Radius != tc.expRadius {
				t.Errorf("Circle radius is (%v) should be (%v)", circle.Radius, tc.expRadius)
			}
		})
	}
}
