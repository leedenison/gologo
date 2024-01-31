package obj

import (
	"fmt"
	"image"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo"
)

func Bitmap(origin mgl32.Vec2, rgba *image.RGBA) *gologo.Object {
	bitmapRenderer, err := gologo.NewBitmapRenderer(rgba)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Bitmap renderer: %v\n", err))
	}

	return &gologo.Object{
		Position: mgl32.Vec3{origin[0], origin[1], 0.0},
		Scale:    1.0,
		Creation: gologo.GetTickTime(),
		ZOrder:   0,
		Renderer: bitmapRenderer,
	}
}
