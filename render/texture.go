package render

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"

	// Bring in png so we support this file format
	_ "image/png"
)

// GLTexture : stores core data on a GL texture
type GLTexture struct {
	ID   uint32
	Size [2]uint32
}

/////////////////////////////////////////////////////////////
// OpenGL Resources
//

// CreateTexture : use the image from the supplied path to create a texture
func CreateTextureImpl(texturePath string) (*GLTexture, error) {
	result, textureExists := glState.Textures[texturePath]
	if !textureExists {
		texture, sizeX, sizeY, err := loadTexture(
			texturePath,
			gl.TEXTURE0)
		if err != nil {
			return nil, err
		}
		result = &GLTexture{
			ID:   texture,
			Size: [2]uint32{sizeX, sizeY},
		}
		glState.Textures[texturePath] = result
	}

	return result, nil
}

func loadTexture(file string, textureUnit uint32) (uint32, uint32, uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to load texture %q: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, 0, 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, 0, 0, fmt.Errorf("unsupported image stride")
	}
	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.Point{}, draw.Src)
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Over)

	return TextureFromRGBA(rgba, textureUnit), uint32(rgba.Rect.Size().X), uint32(rgba.Rect.Size().Y), nil
}

func TextureFromRGBA(rgba *image.RGBA, textureUnit uint32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(textureUnit)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return texture
}
