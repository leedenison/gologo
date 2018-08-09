package gologo

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

func Rectangle(rect Rect, color mgl32.Vec4) *Object {
	originX := (rect[0][0] + rect[1][0]) / 2
	originY := (rect[0][1] + rect[1][1]) / 2

	meshVertices := []float32{
		// Bottom left
		rect[0][0] - originX,
		rect[0][1] - originY,
		0.0,
		0.0,
		1.0,
		// Top right
		rect[1][0] - originX,
		rect[1][1] - originY,
		0.0,
		1.0,
		0.0,
		// Top left
		rect[0][0] - originX,
		rect[1][1] - originY,
		0.0,
		0.0,
		0.0,
		// Bottom left
		rect[0][0] - originX,
		rect[0][1] - originY,
		0.0,
		0.0,
		1.0,
		// Bottom right
		rect[1][0] - originX,
		rect[0][1] - originY,
		0.0,
		1.0,
		1.0,
		// Top right
		rect[1][0] - originX,
		rect[1][1] - originY,
		0.0,
		1.0,
		0.0,
	}

	meshRenderer, err := CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{uniformColor},
		map[int]interface{}{
			uniformColor: color,
		},
		meshVertices)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Rectangle renderer: %v\n", err))
	}

	return &Object{
		Model:    mgl32.Translate3D(originX, originY, 0.0),
		Creation: GetTickTime(),
		ZOrder:   0,
		Renderer: meshRenderer,
	}
}
