package obj

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo"
	"github.com/leedenison/gologo/render"
	"github.com/leedenison/gologo/time"
)

func Rectangle(rect gologo.Rect, color mgl32.Vec4) *gologo.Object {
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

	meshRenderer, err := render.CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{render.UniformColor},
		map[int]interface{}{
			render.UniformColor: color,
		},
		meshVertices)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Rectangle renderer: %v\n", err))
	}

	return &gologo.Object{
		Position: mgl32.Vec3{originX, originY, 0.0},
		Scale:    1.0,
		Creation: time.GetTickTime(),
		ZOrder:   0,
		Renderer: meshRenderer,
	}
}

func Polygon(origin mgl32.Vec2, sides int, radius float32, color mgl32.Vec4) *gologo.Object {
	meshVertices := []float32{}
	angle := 2 * math.Pi / float64(sides)

	for i := 0; i < sides; i++ {
		cos := float32(math.Cos(float64(i) * angle))
		sin := float32(math.Sin(float64(i) * angle))
		meshVertices = append(meshVertices, radius*cos)
		meshVertices = append(meshVertices, radius*sin)
		meshVertices = append(meshVertices, 0.0)
		meshVertices = append(meshVertices, 0.5*cos+0.5)
		meshVertices = append(meshVertices, 0.5*sin+0.5)

		cos = float32(math.Cos(float64(i+1) * angle))
		sin = float32(math.Sin(float64(i+1) * angle))
		meshVertices = append(meshVertices, radius*cos)
		meshVertices = append(meshVertices, radius*sin)
		meshVertices = append(meshVertices, 0.0)
		meshVertices = append(meshVertices, 0.5*cos+0.5)
		meshVertices = append(meshVertices, 0.5*sin+0.5)

		meshVertices = append(meshVertices, 0.0, 0.0, 0.0, 0.5, 0.5)
	}

	meshRenderer, err := render.CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{render.UniformColor},
		map[int]interface{}{
			render.UniformColor: color,
		},
		meshVertices)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Polygon renderer: %v\n", err))
	}

	return &gologo.Object{
		Position: mgl32.Vec3{origin[0], origin[1], 0.0},
		Scale:    1.0,
		Creation: time.GetTickTime(),
		ZOrder:   0,
		Renderer: meshRenderer,
	}
}
