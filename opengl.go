package gologo

import (
	"fmt"
	"image"
	"image/draw"

	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	// Bring in png so we support this file format
	_ "image/png"
)

// GLState : Stores the shaders, textures, and projection
type GLState struct {
	Shaders         map[string]*GLShader
	Textures        map[string]*GLTexture
	NextTextureUnit int32
	Projection      mgl32.Mat4
}

// GLTexture : stores core data on a GL texture
type GLTexture struct {
	ID   uint32
	Size [2]uint32
}

// GLShader : Stores core info for a GL shader
type GLShader struct {
	Program    uint32
	Projection int32
	Model      int32
	Uniforms   map[int]int32
}

// CreateTexture : use the image from the supplied path to create a texture
func CreateTexture(texturePath string) (*GLTexture, error) {
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
		return 0, 0, 0, fmt.Errorf("Failed to load texture %q: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, 0, 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, 0, 0, fmt.Errorf("Unsupported image stride")
	}
	draw.Draw(rgba, rgba.Bounds(), image.Transparent, image.ZP, draw.Src)
	draw.Draw(rgba, rgba.Bounds(), img, image.ZP, draw.Over)

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

	return texture, uint32(rgba.Rect.Size().X), uint32(rgba.Rect.Size().Y), nil
}

// CreateShaderProgram : compiles the vertex and fragment shaders, create the program
// attach the shaders, link the program and store the program with it's uniforms for later use
func CreateShaderProgram(vertexShader string, fragmentShader string) (*GLShader, error) {
	programKey := vertexShader + "," + fragmentShader
	program, programExists := glState.Shaders[programKey]
	if !programExists {
		var err error
		program, err = loadProgram(shaders[vertexShader], shaders[fragmentShader])
		if err != nil {
			return nil, err
		}
		glState.Shaders[programKey] = program
	}

	return program, nil
}

func loadProgram(vertexShaderSource, fragmentShaderSource string) (*GLShader, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return &GLShader{
		Program:  program,
		Uniforms: map[int]int32{},
	}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", log, source)
	}

	return shader, nil
}
