package gologo

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

// Need to decide on corner cases like all 0s
var simplePrimTests = []struct {
	botX, topX, botY, topY, radius, posX, posY float32
	contExp, overlExp                          bool
}{
	{0, 400, 0, 400, 50, 200, 200, true, true},
	{0, 400, 0, 400, 50, 380, 380, false, true},
	{0, 400, 0, 400, 50, 425, 425, false, true},
	{0, 400, 0, 400, 50, 475, 475, false, false},
	{0, 400, 0, 400, 20, 425, 425, false, false},
	{0, 400, 0, 400, 100, 500, 400, false, true},
	{0, 400, 0, 400, 100, 501, 400, false, false},
}

// TestSimplePrimitive : Test basic primitive creation and
// functionality
func TestSimplePrimitive(t *testing.T) {
	var x, y float32
	var rect Rect
	var position mgl32.Vec3
	var circle Circle
	var isContained, overlaps bool

	obj := CreateObject(position)

	for _, tt := range simplePrimTests {
		obj.SetPosition(tt.posX, tt.posY)
		circle.Radius = tt.radius
		x, y = obj.GetPosition()
		rect = Rect{{tt.botX, tt.botY}, {tt.topX, tt.topY}}

		isContained = circle.IsContainedInRect(*obj, rect)
		if isContained != tt.contExp {
			t.Errorf("IsContInRect %v (exp: %v) with x (%v), y (%v), rad (%v) in x (%v-%v) and y (%v-%v)",
				isContained, tt.contExp, x, y, circle.Radius, tt.botX, tt.topX, tt.botY, tt.topY)
		}

		overlaps = circle.OverlapsWithRect(*obj, rect)
		if overlaps != tt.overlExp {
			t.Errorf("overlapWithRect %v (exp: %v) with x (%v), y (%v), rad (%v) in x (%v-%v) and y (%v-%v)",
				overlaps, tt.overlExp, x, y, circle.Radius, tt.botX, tt.topX, tt.botY, tt.topY)
		}
	}
}
