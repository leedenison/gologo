package opengl

import (
	"fmt"
	"image"
	"image/draw"

	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/render"

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

func InitOpenGL() error {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		log.Error.Println("gl.Init failed:", err)
		return err
	}

	CreateMeshRenderer = CreateMeshRendererImpl
	CreateTexture = CreateTextureImpl

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Trace.Println("OpenGL version:", version)

	// Configure global settings
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	return nil
}

func UpdateProjection(width float32, height float32) {
	glState.Projection = mgl32.Ortho2D(0, width, 0, height)
}

func ClearBackBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
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

/////////////////////////////////////////////////////////////
// MeshRenderer
//

type MeshRenderer struct {
	Shader       *GLShader
	Mesh         uint32
	Uniforms     map[int]interface{}
	MeshVertices []float32
	VertexCount  int32
}

func (r *MeshRenderer) Render(model mgl32.Mat4) {
	r.RenderAt(model, map[int]interface{}{})
}

func (r *MeshRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	gl.UseProgram(r.Shader.Program)
	gl.UniformMatrix4fv(r.Shader.Model, 1, false, &model[0])
	gl.UniformMatrix4fv(r.Shader.Projection, 1, false, &glState.Projection[0])

	r.bindCustomUniforms(r.Shader, custom)

	gl.BindVertexArray(r.Mesh)

	gl.DrawArrays(gl.TRIANGLES, 0, r.VertexCount)
}

// Binds statically defined and custom uniforms.  Custom uniforms take
// precedence over statically defined uniforms.
func (r *MeshRenderer) bindCustomUniforms(shader *GLShader, custom map[int]interface{}) {
	glState.NextTextureUnit = 0

	for location, value := range r.Uniforms {
		if _, exists := custom[location]; !exists {
			r.bindCustomUniform(shader, location, value)
		}
	}

	for location, value := range custom {
		r.bindCustomUniform(shader, location, value)
	}
}

func (r *MeshRenderer) bindCustomUniform(
	shader *GLShader, location int, value interface{}) {
	switch tValue := value.(type) {
	case *GLTexture:
		gl.ActiveTexture(gl.TEXTURE0 + uint32(glState.NextTextureUnit))
		gl.Uniform1i(shader.Uniforms[location], glState.NextTextureUnit)
		gl.BindTexture(gl.TEXTURE_2D, tValue.ID)
		glState.NextTextureUnit++
	case int32:
		gl.Uniform1i(shader.Uniforms[location], tValue)
	case float32:
		gl.Uniform1f(shader.Uniforms[location], tValue)
	case mgl32.Vec4:
		gl.Uniform4fv(shader.Uniforms[location], 1, &tValue[0])
	default:
		panic(fmt.Sprintf("Unhandled uniform(%v) value type: %t\n", location, value))
	}
}

func (r *MeshRenderer) DebugRender(model mgl32.Mat4) {
	//r.DebugRenderAt(model, map[int]interface{}{})
}

func (r *MeshRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	log.Trace.Printf("MeshRenderer: Model matrix:\n%v\n", model)
	log.Trace.Printf("MeshRenderer: Mesh vertices:\n%v\n", r.MeshVertices)

	rendered := []mgl32.Vec4{}
	for i := 0; i < len(r.MeshVertices); i += GlMeshStride {
		rv := mgl32.Vec4{
			r.MeshVertices[i],
			r.MeshVertices[i+1],
			r.MeshVertices[i+2],
			1,
		}

		rv = model.Mul4x1(rv)

		rendered = append(rendered, rv)
	}

	log.Trace.Printf("MeshRenderer: Rendered vertices:\n%v\n", rendered)
}

func (r *MeshRenderer) Animate(model mgl32.Mat4) {}

// Clone : Clones a MeshRenderer.  The shader program and mesh are
// retained as shared instances.  The uniform values are shallow copied.
func (r *MeshRenderer) Clone() render.Renderer {
	uniforms := make(map[int]interface{})
	for k, v := range r.Uniforms {
		uniforms[k] = v
	}

	return &MeshRenderer{
		Shader:       r.Shader,
		Mesh:         r.Mesh,
		Uniforms:     uniforms,
		MeshVertices: r.MeshVertices,
		VertexCount:  r.VertexCount,
	}
}

// CreateMeshRenderer : Creates a MeshRenderer.  uniforms specifies all uniform variable locations that
// should be bound in the shader program, including both uniforms with a statically
// defined value and those supplied in each call to RenderAt.  uniformValues
// specifies static values for some or all of the uniforms.
func CreateMeshRendererImpl(
	vertexShader string,
	fragmentShader string,
	uniforms []int,
	uniformValues map[int]interface{},
	meshVertices []float32) (*MeshRenderer, error) {
	shader, err := CreateShaderProgram(vertexShader, fragmentShader)
	if err != nil {
		return nil, err
	}

	gl.UseProgram(shader.Program)
	gl.BindFragDataLocation(shader.Program, 0, fragLocOutputColor)

	shader.Projection = gl.GetUniformLocation(shader.Program, shaderUniformLocProjection)
	shader.Model = gl.GetUniformLocation(shader.Program, shaderUniformLocModel)

	for _, uniform := range uniforms {
		shader.Uniforms[uniform] = gl.GetUniformLocation(shader.Program, shaderUniforms[uniform])
	}

	mesh := createMeshBuffer(shader.Program, meshVertices)

	return &MeshRenderer{
		Shader:       shader,
		Mesh:         mesh,
		Uniforms:     uniformValues,
		MeshVertices: meshVertices,
		VertexCount:  int32(len(meshVertices) / GlMeshStride),
	}, nil
}

func createMeshBuffer(shader uint32, vertices []float32) uint32 {
	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(vertices)*4,
		gl.Ptr(vertices),
		gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(shader, attribLocVertex))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, glMeshStrideBytes,
		gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(shader, attribLocVertexTexCoord))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, glMeshStrideBytes,
		gl.PtrOffset(3*4))

	return vao
}
