package render

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/time"
)

// GLState : Stores the shaders, textures, and projection
type GLState struct {
	Shaders         map[string]*GLShader
	Textures        map[string]*GLTexture
	NextTextureUnit int32
	Projection      mgl32.Mat4
}

type Renderer interface {
	Render(model mgl32.Mat4)
	RenderAt(model mgl32.Mat4, uniforms map[int]interface{})
	DebugRender(model mgl32.Mat4)
	DebugRenderAt(model mgl32.Mat4, uniforms map[int]interface{})
	Animate(model mgl32.Mat4)
	Clone() Renderer
}

/////////////////////////////////////////////////////////////
// Rendering globals
//

const (
	GlMeshStride      = 5
	glMeshStrideBytes = GlMeshStride * float32SizeBytes
)

var CreateMeshRenderer func(
	vertexShader string,
	fragmentShader string,
	uniforms []int,
	uniformValues map[int]interface{},
	meshVertices []float32) (*MeshRenderer, error)

var CreateTexture func(texturePath string) (*GLTexture, error)

var glState = &GLState{
	Shaders:         map[string]*GLShader{},
	Textures:        map[string]*GLTexture{},
	NextTextureUnit: gl.TEXTURE0,
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

func ClearBackBuffer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func Set2DProjection(width float32, height float32) {
	glState.Projection = mgl32.Ortho2D(0, width, 0, height)
}

func Colors(gradients []float64, count int) []color.RGBA {
	stride := 4
	red := 1
	green := 2
	blue := 3

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

/////////////////////////////////////////////////////////////
// BitmapRenderer
//

type BitmapRenderer struct {
	MeshRenderer *MeshRenderer
	Buffer       *image.RGBA
}

func (r *BitmapRenderer) Render(model mgl32.Mat4) {
	r.RenderAt(model, map[int]interface{}{})
}

func (r *BitmapRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	r.MeshRenderer.RenderAt(model, custom)
}

func (r *BitmapRenderer) DebugRender(model mgl32.Mat4) {
	r.DebugRenderAt(model, map[int]interface{}{})
}

func (r *BitmapRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	r.MeshRenderer.DebugRenderAt(model, custom)
}

func (r *BitmapRenderer) Animate(model mgl32.Mat4) {
	texture, ok := r.MeshRenderer.Uniforms[UniformTexture].(*GLTexture)
	if !ok {
		return
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)
	gl.TexSubImage2D(
		gl.TEXTURE_2D,
		0,
		0,
		0,
		int32(r.Buffer.Rect.Size().X),
		int32(r.Buffer.Rect.Size().Y),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(r.Buffer.Pix))
}

func (r *BitmapRenderer) Clone() Renderer {
	return &BitmapRenderer{
		MeshRenderer: r.MeshRenderer,
		Buffer:       r.Buffer,
	}
}

func NewBitmapRenderer(rgba *image.RGBA) (*BitmapRenderer, error) {
	meshVertices := []float32{
		// Bottom left
		-1.0, -1.0, 0.0, 0.0, 1.0,
		// Top right
		1.0, 1.0, 0.0, 1.0, 0.0,
		// Top left
		-1.0, 1.0, 0.0, 0.0, 0.0,
		// Bottom left
		-1.0, -1.0, 0.0, 0.0, 1.0,
		// Bottom right
		1.0, -1.0, 0.0, 1.0, 1.0,
		// Top right
		1.0, 1.0, 0.0, 1.0, 0.0,
	}

	texture := &GLTexture{
		ID:   TextureFromRGBA(rgba, gl.TEXTURE0),
		Size: [2]uint32{uint32(rgba.Rect.Size().X), uint32(rgba.Rect.Size().Y)},
	}

	meshRenderer, err := CreateMeshRenderer(
		"FULLSCREEN_VERTEX_SHADER",
		"TEXTURE_FRAGMENT_SHADER",
		[]int{UniformTexture},
		map[int]interface{}{
			UniformTexture: texture,
		},
		meshVertices)
	if err != nil {
		return nil, err
	}

	return &BitmapRenderer{
		MeshRenderer: meshRenderer,
		Buffer:       rgba,
	}, nil
}

/////////////////////////////////////////////////////////////
// TextRenderer
//

type TextRenderer struct {
	MeshRenderers map[byte]Renderer
	CharWidths    map[byte]float32
	CharSpacer    float32
	Text          []byte
	Transforms    []mgl32.Mat4
}

func (r *TextRenderer) Render(model mgl32.Mat4) {
	r.RenderAt(model, map[int]interface{}{})
}

func (r *TextRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	if len(r.Transforms) != len(r.Text) {
		r.InitTransforms()
	}

	for i := 0; i < len(r.Text); i++ {
		renderer, ok := r.MeshRenderers[r.Text[i]]

		if ok {
			renderer.RenderAt(
				model.Mul4(r.Transforms[i]),
				map[int]interface{}{})
		}
	}
}

func (r *TextRenderer) InitTransforms() {
	var translate float32

	if len(r.Text) == 0 {
		return
	}

	r.Transforms = make([]mgl32.Mat4, len(r.Text))
	for count := 0; count < len(r.Text); count++ {
		width, ok := r.CharWidths[r.Text[count]]
		if ok {
			r.Transforms[count] = mgl32.Translate3D(
				translate+(width/2), 0, 0)
			translate += width + r.CharSpacer
		}
	}

	for count := 0; count < len(r.Text); count++ {
		r.Transforms[count] = mgl32.Translate3D(-translate/2, 0, 0).
			Mul4(r.Transforms[count])
	}
}

func (r *TextRenderer) DebugRender(model mgl32.Mat4) {
	r.DebugRenderAt(model, map[int]interface{}{})
}

func (r *TextRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	if len(r.Transforms) != len(r.Text) {
		r.InitTransforms()
	}

	for i := 0; i < len(r.Text); i++ {
		renderer, ok := r.MeshRenderers[r.Text[i]]

		if ok {
			renderer.DebugRenderAt(
				model.Mul4(r.Transforms[i]),
				map[int]interface{}{})
		}
	}
}

func (r *TextRenderer) Animate(model mgl32.Mat4) {}

func (r *TextRenderer) Clone() Renderer {
	return &TextRenderer{
		MeshRenderers: r.MeshRenderers,
		CharSpacer:    r.CharSpacer,
	}
}

/////////////////////////////////////////////////////////////
// ExplosionRenderer
//

type Particle struct {
	Model    mgl32.Mat4
	Velocity mgl32.Mat4
	Age      float64
	Renderer Renderer
}

type ExplosionRenderer struct {
	Renderers     []Renderer
	ParticleCount int
	Particles     []*Particle
	MaxAge        float64
}

func (r *ExplosionRenderer) Render(model mgl32.Mat4) {
	r.RenderAt(model, map[int]interface{}{})
}

func (r *ExplosionRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	for i := 0; i < len(r.Particles); i++ {
		age := float64(time.GetTickTime()) - r.Particles[i].Age
		r.Particles[i].Renderer.RenderAt(
			model.Mul4(r.Particles[i].Model),
			map[int]interface{}{
				UniformAlpha: 1.0 - age/r.MaxAge,
			})
	}
}

func (r *ExplosionRenderer) DebugRender(model mgl32.Mat4) {
	r.DebugRenderAt(model, map[int]interface{}{})
}

func (r *ExplosionRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	log.Trace.Printf("ExplosionRenderer: Model matrix:\n%v\n", model)
}

func (r *ExplosionRenderer) Animate(model mgl32.Mat4) {
	if r.Particles == nil {
		r.Particles = make([]*Particle, r.ParticleCount)
		for i := 0; i < r.ParticleCount; i++ {
			r.Particles[i] = r.createRandomParticle()
		}
	}

	for j := 0; j < len(r.Particles); j++ {
		age := float64(time.GetTickTime()) - r.Particles[j].Age
		if age >= r.MaxAge {
			if len(r.Particles) > 1 {
				r.Particles = append(r.Particles[:j], r.Particles[j+1:]...)
			} else {
				r.Particles = r.Particles[0:0]
			}
			j--
		} else {
			r.Particles[j].Model = r.Particles[j].Velocity.Mul4(r.Particles[j].Model)
		}
	}
}

func (r *ExplosionRenderer) createRandomParticle() *Particle {
	renderer := rand.Intn(len(r.Renderers))
	velocity := mgl32.HomogRotate3DZ(rand.Float32() * 0.01)
	velocity = velocity.Mul4(mgl32.Translate3D(rand.Float32(), rand.Float32(), 0.0))

	return &Particle{
		Renderer: r.Renderers[renderer],
		Velocity: velocity,
		Age:      float64(time.GetTickTime()),
		Model:    mgl32.Ident4(),
	}
}

func (r *ExplosionRenderer) Clone() Renderer {
	return &ExplosionRenderer{
		Renderers: r.Renderers,
		MaxAge:    r.MaxAge,
	}
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
	shader *GLShader, location int, value interface{},
) {
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
	// r.DebugRenderAt(model, map[int]interface{}{})
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
func (r *MeshRenderer) Clone() Renderer {
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
	meshVertices []float32,
) (*MeshRenderer, error) {
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
