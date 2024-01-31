package gologo

import (
	"fmt"
	"image"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo/opengl"
	"github.com/leedenison/gologo/time"
)

func Bitmap(origin mgl32.Vec2, rgba *image.RGBA) *Object {
	bitmap := InvisibleBitmap(origin, rgba)
	TagRender(bitmap)
	return bitmap
}

func InvisibleBitmap(origin mgl32.Vec2, rgba *image.RGBA) *Object {
	bitmapRenderer, err := NewBitmapRenderer(rgba)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Bitmap renderer: %v\n", err))
	}

	return &Object{
		Position: mgl32.Vec3{origin[0], origin[1], 0.0},
		Scale:    1.0,
		Creation: time.GetTickTime(),
		ZOrder:   0,
		Renderer: bitmapRenderer,
	}
}

func Rectangle(rect Rect, color mgl32.Vec4) *Object {
	rectangle := InvisibleRectangle(rect, color)
	TagRender(rectangle)
	return rectangle
}

func InvisibleRectangle(rect Rect, color mgl32.Vec4) *Object {
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

	meshRenderer, err := opengl.CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{opengl.UniformColor},
		map[int]interface{}{
			opengl.UniformColor: color,
		},
		meshVertices)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Rectangle renderer: %v\n", err))
	}

	return &Object{
		Position: mgl32.Vec3{originX, originY, 0.0},
		Scale:    1.0,
		Creation: time.GetTickTime(),
		ZOrder:   0,
		Renderer: meshRenderer,
	}
}

func Polygon(origin mgl32.Vec2, sides int, radius float32, color mgl32.Vec4) *Object {
	polygon := InvisiblePolygon(origin, sides, radius, color)
	TagRender(polygon)
	return polygon
}

func InvisiblePolygon(origin mgl32.Vec2, sides int, radius float32, color mgl32.Vec4) *Object {
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

	meshRenderer, err := opengl.CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{opengl.UniformColor},
		map[int]interface{}{
			opengl.UniformColor: color,
		},
		meshVertices)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Polygon renderer: %v\n", err))
	}

	return &Object{
		Position: mgl32.Vec3{origin[0], origin[1], 0.0},
		Scale:    1.0,
		Creation: time.GetTickTime(),
		ZOrder:   0,
		Renderer: meshRenderer,
	}
}
