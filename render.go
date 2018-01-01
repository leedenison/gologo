package gologo

import (
    "fmt"
    "math/rand"
    "sort"
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/mathgl/mgl32"
)

type ScreenDirection int

const (
    SCREEN_UP ScreenDirection = iota
    SCREEN_DOWN
    SCREEN_LEFT
    SCREEN_RIGHT
)

type Rect [2][2]float32

type Renderer interface {
    Render(object *Object)
    DebugRender(object *Object)
    Animate(object *Object)
    Clone() Renderer
}

func InitRender() error {
    // Initialize Glow
    if err := gl.Init(); err != nil {
        Error.Println("gl.Init failed:", err)
        return err
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    Trace.Println("OpenGL version:", version)

    // Configure global settings
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)

    return nil
}

func UpdateWindowProjection() {
    glState.Projection = mgl32.Ortho2D(
        0, float32(windowState.Width), 0, float32(windowState.Height))
}

func ClearBackBuffer() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)    
}

func Render() {
    // Sort the objects by the rendering zorder
    sort.Sort(ByZOrder(objects))
    for _, object := range objects {
        object.Renderer.Animate(object)
        //object.ObjectType.Renderer.DebugRender(object)
        object.Renderer.Render(object)
    }
    windowState.Main.SwapBuffers()
}

/////////////////////////////////////////////////////////////
// MeshRenderer
//

type MeshRenderer struct {
    Shader *GLShader
    Mesh uint32
    Uniforms map[int]interface{}
    MeshVertices []float32
    VertexCount int32
}

func (r *MeshRenderer) Render(object *Object) {
    r.RenderAt(object.Model, map[int]interface{}{})
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
    default:
        panic(fmt.Sprintf("Unhandled uniform(%v) value type: %t\n", location, value))
    }
}

func (r *MeshRenderer) DebugRender(object *Object) {
    //r.DebugRenderAt(object.Model, map[int]interface{}{})
}

func (r *MeshRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
    Trace.Printf("MeshRenderer: Model matrix:\n%v\n", model)
    Trace.Printf("MeshRenderer: Mesh vertices:\n%v\n", r.MeshVertices)

    rendered := []mgl32.Vec4 {}
    for i := 0; i < len(r.MeshVertices); i += GL_MESH_STRIDE {
        rv := mgl32.Vec4 {
            r.MeshVertices[i],
            r.MeshVertices[i + 1],
            r.MeshVertices[i + 2],
            1,
        }

        rv = model.Mul4x1(rv)

        rendered = append(rendered, rv)
    }

    Trace.Printf("MeshRenderer: Rendered vertices:\n%v\n", rendered)
}

func (r *MeshRenderer) Animate(object *Object) {}

// Clones a MeshRenderer.  The shader program and mesh are retained as shared
// instances.  The uniform values are shallow copied.
func (r *MeshRenderer) Clone() Renderer {
    uniforms := make(map[int]interface{})
    for k, v := range r.Uniforms {
        uniforms[k] = v
    }

    return &MeshRenderer {
        Shader: r.Shader,
        Mesh: r.Mesh,
        Uniforms: uniforms,
        MeshVertices: r.MeshVertices,
        VertexCount: r.VertexCount,
    }
}

// Creates a MeshRenderer.  uniforms specifies all uniform variable locations that
// should be bound in the shader program, including both uniforms with a statically
// defined value and those supplied in each call to RenderAt.  uniformValues 
// specifies static values for some or all of the uniforms.
func CreateMeshRenderer(
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
    gl.BindFragDataLocation(shader.Program, 0, FRAG_LOC_OUTPUT_COLOR)

    shader.Projection = gl.GetUniformLocation(shader.Program, UNIFORM_LOC_PROJECTION)
    shader.Model = gl.GetUniformLocation(shader.Program, UNIFORM_LOC_MODEL)


    for _, uniform := range uniforms {
        shader.Uniforms[uniform] = gl.GetUniformLocation(shader.Program, UNIFORMS[uniform])
    }

    mesh := createMeshBuffer(shader.Program, meshVertices)

    return &MeshRenderer {
        Shader: shader,
        Mesh: mesh,
        Uniforms: uniformValues,
        MeshVertices: meshVertices,
        VertexCount: int32(len(meshVertices) / GL_MESH_STRIDE),
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

    vertAttrib := uint32(gl.GetAttribLocation(shader, ATTRIB_LOC_VERTEX))
    gl.EnableVertexAttribArray(vertAttrib)
    gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, GL_MESH_STRIDE_BYTES,
        gl.PtrOffset(0))

    texCoordAttrib := uint32(gl.GetAttribLocation(shader, ATTRIB_LOC_VERTEX_TEX_COORD))
    gl.EnableVertexAttribArray(texCoordAttrib)
    gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, GL_MESH_STRIDE_BYTES,
        gl.PtrOffset(3*4))

    return vao
}

/////////////////////////////////////////////////////////////
// TextureRenderer
//

func CreateTextureRenderer(
        vertexShader string,
        fragmentShader string,
        texture *GLTexture,
        uniforms []int,
        uniformValues map[int]interface{},
        meshVertices []float32) (*MeshRenderer, error) {
    if uniformValues == nil {
        uniformValues = map[int]interface{}{}
    }

    uniformValues[UNIFORM_TEXTURE] = texture

    if !containsInt(uniforms, UNIFORM_TEXTURE) {
        uniforms = append(uniforms, UNIFORM_TEXTURE)
    }

    return CreateMeshRenderer(
        vertexShader, fragmentShader, uniforms, uniformValues, meshVertices)
}

/////////////////////////////////////////////////////////////
// TextRenderer
//

type TextRenderer struct {
    MeshRenderers map[byte]*MeshRenderer
    CharWidths map[byte]float32
    CharSpacer float32
    Text []byte
    Transforms []mgl32.Mat4
}

func (r *TextRenderer) Render(object *Object) {
    if len(r.Transforms) != len(r.Text) {
        r.InitTransforms()
    }

    for i := 0; i < len(r.Text); i++ {
        renderer, ok := r.MeshRenderers[r.Text[i]]

        if ok {
            renderer.RenderAt(
                object.Model.Mul4(r.Transforms[i]),
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
                translate + (width / 2), 0, 0)
            translate += width + r.CharSpacer
        }
    }

    for count := 0; count < len(r.Text); count++ {
        r.Transforms[count] = mgl32.Translate3D(-translate / 2, 0, 0).
            Mul4(r.Transforms[count])
    }
}

func (r *TextRenderer) DebugRender(object *Object) {
    if len(r.Transforms) != len(r.Text) {
        r.InitTransforms()
    }

    for i := 0; i < len(r.Text); i++ {
        renderer, ok := r.MeshRenderers[r.Text[i]]

        if ok {
            renderer.DebugRenderAt(
                object.Model.Mul4(r.Transforms[i]),
                map[int]interface{} {})
        }
    }
}

func (r *TextRenderer) Animate(object *Object) {}

func (r *TextRenderer) Clone() Renderer {
    return &TextRenderer {
        MeshRenderers: r.MeshRenderers,
        CharSpacer: r.CharSpacer, 
    }
}

func CalcMeshFromChar(
    textureSizeX uint32,
    textureSizeY uint32,
    textureRect [][2]float32,
    charRect [][2]float32) []float32 {

    textureWidth := textureRect[1][0] - textureRect[0][0]
    textureHeight := textureRect[1][1] - textureRect[0][1]
    deltaRight := textureRect[1][0] - charRect[1][0]
    deltaLeft := charRect[0][0] - textureRect[0][0]
    widthDelta := deltaRight - deltaLeft
    deltaTop := charRect[0][1] - textureRect[0][1]
    deltaBottom := textureRect[1][1] - charRect[1][1]
    heightDelta := deltaTop - deltaBottom

    return []float32 {
        // Bottom left
        (-float32(textureWidth) + widthDelta) / 2,
        (-float32(textureHeight) + heightDelta) / 2,
        0.0,
        textureRect[0][0] / float32(textureSizeX),
        textureRect[1][1] / float32(textureSizeY),
        // Top right
        (float32(textureWidth) + widthDelta) / 2,
        (float32(textureHeight) + heightDelta) / 2,
        0.0,
        textureRect[1][0] / float32(textureSizeX),
        textureRect[0][1] / float32(textureSizeY),
        // Top left
        (-float32(textureWidth) + widthDelta) / 2,
        (float32(textureHeight) + heightDelta) / 2,
        0.0,
        textureRect[0][0] / float32(textureSizeX),
        textureRect[0][1] / float32(textureSizeY),
        // Bottom left
        (-float32(textureWidth) + widthDelta) / 2,
        (-float32(textureHeight) + heightDelta) / 2,
        0.0,
        textureRect[0][0] / float32(textureSizeX),
        textureRect[1][1] / float32(textureSizeY),
        // Bottom right
        (float32(textureWidth) + widthDelta) / 2,
        (-float32(textureHeight) + heightDelta) / 2,
        0.0,
        textureRect[1][0] / float32(textureSizeX),
        textureRect[1][1] / float32(textureSizeY),
        // Top right
        (float32(textureWidth) + widthDelta) / 2,
        (float32(textureHeight) + heightDelta) / 2,
        0.0,
        textureRect[1][0] / float32(textureSizeX),
        textureRect[0][1] / float32(textureSizeY),
    }
}

/////////////////////////////////////////////////////////////
// ExplosionRenderer
//

type Particle struct {
    Model mgl32.Mat4
    Velocity mgl32.Mat4
    Renderer *MeshRenderer
}

type ExplosionRenderer struct {
    MeshRenderers []*MeshRenderer
    ParticleCount int
    Particles []*Particle
    MaxAge float32
}

func (r *ExplosionRenderer) Render(object *Object) {
    age := GetTickTime() - object.Creation
    for i := 0; i < len(r.Particles); i++ {
        r.Particles[i].Renderer.RenderAt(
            object.Model.Mul4(r.Particles[i].Model),
            map[int]interface{} {
                UNIFORM_ALPHA: 1.0 - float32(age) / float32(r.MaxAge),
            })
    }
}

func (r *ExplosionRenderer) DebugRender(object *Object) {
    Trace.Printf("ExplosionRenderer: Model matrix:\n%v\n", object.Model)
}

func (r *ExplosionRenderer) Animate(object *Object) {
    age := float32(GetTickTime() - object.Creation)
    if r.Particles == nil {
        r.Particles = make([]*Particle, r.ParticleCount)
        for i:= 0; i < r.ParticleCount; i++ {
            r.Particles[i] = r.createRandomParticle()
        }
    }

    for j := 0; j < len(r.Particles); j++ {
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
    renderer := rand.Intn(len(r.MeshRenderers))
    velocity := mgl32.HomogRotate3DZ(rand.Float32() * 0.01)
    velocity = velocity.Mul4(mgl32.Translate3D(rand.Float32(), rand.Float32(), 0.0))

    return &Particle {
        Renderer: r.MeshRenderers[renderer],
        Velocity: velocity,
        Model: mgl32.Ident4(),
    }
}

func (r *ExplosionRenderer) Clone() Renderer {
    return &ExplosionRenderer {
        MeshRenderers: r.MeshRenderers,
        MaxAge: r.MaxAge,
    }
}