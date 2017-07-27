package gologo

import (
    "fmt"
    "math/rand"
    "github.com/pkg/errors"
    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/mathgl/mgl32"
)

type Renderer interface {
    Render(object *Object)
    DebugRender(object *Object)
    Animate(object *Object)
}

type GLTexture struct {
    ID uint32
    TextureUnit uint32
}

/////////////////////////////////////////////////////////////
// GLMeshRenderer
//

type GLMeshRenderer struct {
    Shader *GLShader
    Texture *GLTexture
    Mesh uint32
    MeshVertices []float32
    VertexCount int32
}

func (r *GLMeshRenderer) Render(object *Object) {
    r.RenderAt(object.Model, map[int]interface{}{})
}

func (r *GLMeshRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
    gl.UseProgram(r.Shader.Program)
    gl.UniformMatrix4fv(r.Shader.Model, 1, false, &model[0])
    gl.UniformMatrix4fv(r.Shader.Projection, 1, false, &glState.Projection[0])

    bindCustomUniforms(r.Shader, custom)

    gl.BindVertexArray(r.Mesh)

    if r.Texture != nil {
        gl.ActiveTexture(gl.TEXTURE0)
        gl.Uniform1i(r.Shader.Texture, 0)
        gl.BindTexture(gl.TEXTURE_2D, r.Texture.ID)
    }

    gl.DrawArrays(gl.TRIANGLES, 0, r.VertexCount)
}

func (r *GLMeshRenderer) DebugRender(object *Object) {
    Trace.Printf("GLMeshRenderer: Model matrix:\n%v\n", object.Model)
    Trace.Printf("GLMeshRenderer: Mesh vertices:\n%v\n", r.MeshVertices)

    rendered := []mgl32.Vec4 {}
    for i := 0; i < len(r.MeshVertices); i += GL_MESH_STRIDE {
        rv := mgl32.Vec4 {
            r.MeshVertices[i],
            r.MeshVertices[i + 1],
            r.MeshVertices[i + 2],
            1,
        }

        rv = object.Model.Mul4x1(rv)

        rendered = append(rendered, rv)
    }

    Trace.Printf("GLMeshRenderer: Rendered vertices:\n%v\n", rendered)
}

func (r *GLMeshRenderer) Animate(object *Object) {}

func bindCustomUniforms(shader *GLShader, custom map[int]interface{}) {
    for location, value := range custom {
        switch tValue := value.(type) {
        case int32:
            gl.Uniform1i(shader.Uniforms[location], tValue)
        case float32:
            gl.Uniform1f(shader.Uniforms[location], tValue)
        default:
            panic(fmt.Sprintf("Unhandled uniform(%v) value type: %t\n", location, value))
        }
    }
}

func InitSpriteMeshRenderer(config *SpriteMeshRendererConfig) (Renderer, error) {
    if config.Texture == "" {
        return nil, errors.New("Missing required field 'Texture'.")
    }

    texture, sizeX, sizeY, err := InitTexture(config.Texture)
    if err != nil {
        return nil, err
    }

    meshVertices := []float32 {
        // Bottom left
        (-float32(sizeX) / 2 - float32(config.TextureOrigin[0])) * config.MeshScaling,
        (-float32(sizeY) / 2 + float32(config.TextureOrigin[1])) * config.MeshScaling,
        0.0,
        0.0,
        1.0,
        // Top right
        (float32(sizeX) / 2 - float32(config.TextureOrigin[0])) * config.MeshScaling,
        (float32(sizeY) / 2 + float32(config.TextureOrigin[1])) * config.MeshScaling,
        0.0,
        1.0,
        0.0,
        // Top left
        (-float32(sizeX) / 2 - float32(config.TextureOrigin[0])) * config.MeshScaling,
        (float32(sizeY) / 2 + float32(config.TextureOrigin[1])) * config.MeshScaling,
        0.0,
        0.0,
        0.0,
        // Bottom left
        (-float32(sizeX) / 2 - float32(config.TextureOrigin[0])) * config.MeshScaling,
        (-float32(sizeY) / 2 + float32(config.TextureOrigin[1])) * config.MeshScaling,
        0.0,
        0.0,
        1.0,
        // Bottom right
        (float32(sizeX) / 2 - float32(config.TextureOrigin[0])) * config.MeshScaling,
        (-float32(sizeY) / 2 + float32(config.TextureOrigin[1])) * config.MeshScaling,
        0.0,
        1.0,
        1.0,
        // Top right
        (float32(sizeX) / 2 - float32(config.TextureOrigin[0])) * config.MeshScaling,
        (float32(sizeY) / 2 + float32(config.TextureOrigin[1])) * config.MeshScaling,
        0.0,
        1.0,
        0.0,
    }

    return InitMeshRenderer(
        config.VertexShader,
        config.FragmentShader,
        texture,
        []int {},
        meshVertices)
}

func InitGLMeshRenderer(config *GLMeshRendererConfig) (Renderer, error) {
    var texture *GLTexture
    var err error
    if config.Texture != "" {
        texture, _, _, err = InitTexture(config.Texture)
        if err != nil {
            return nil, err
        }
    }

    return InitMeshRenderer(
        config.VertexShader,
        config.FragmentShader,
        texture,
        []int {},
        config.MeshVertices)
}

func InitMeshRenderer(
        vertexShader string,
        fragmentShader string,
        texture *GLTexture,
        uniforms []int,
        meshVertices []float32) (*GLMeshRenderer, error) {
    err := ValidateMeshRenderConfig(vertexShader, fragmentShader, meshVertices)
    if err != nil {
        return nil, err
    }

    shader, err := InitShaderProgram(vertexShader, fragmentShader)
    if err != nil {
        return nil, err
    }

    gl.UseProgram(shader.Program)
    gl.BindFragDataLocation(shader.Program, 0, FRAG_LOC_OUTPUT_COLOR)

    shader.Projection = gl.GetUniformLocation(shader.Program, UNIFORM_LOC_PROJECTION)
    shader.Model = gl.GetUniformLocation(shader.Program, UNIFORM_LOC_MODEL)

    if texture != nil {
        shader.Texture = gl.GetUniformLocation(shader.Program, UNIFORM_LOC_TEXTURE)
    }

    for _, uniform := range uniforms {
        shader.Uniforms[uniform] = gl.GetUniformLocation(shader.Program, UNIFORMS[uniform])
    }

    mesh := InitObjectMesh(shader.Program, meshVertices)

    return &GLMeshRenderer {
        Shader: shader,
        Texture: texture,
        Mesh: mesh,
        MeshVertices: meshVertices,
        VertexCount: int32(len(meshVertices) / GL_MESH_STRIDE),
    }, nil
}

func InitObjectMesh(shader uint32, vertices []float32) uint32 {
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

func ValidateMeshRenderConfig(
        vertexShader string,
        fragmentShader string,
        meshVertices []float32) error {
    if vertexShader == "" {
        return errors.New("Missing required field: 'VertexShader'")
    } else if _, ok := SHADERS[vertexShader]; !ok {
        return errors.Errorf("Unknown 'VertexShader': %v", vertexShader)
    }

    if fragmentShader == "" {
        return errors.New("Missing required field: 'FragmentShader'")
    } else if _, ok := SHADERS[fragmentShader]; !ok {
        return errors.Errorf("Unknown 'FragmentShader': %v", fragmentShader)
    }

    if len(meshVertices) == 0 {
        return errors.New("Missing required field: 'MeshVertices'")
    }

    return nil
}

/////////////////////////////////////////////////////////////
// ExplosionRenderer
//

type Particle struct {
    Model mgl32.Mat4
    Velocity mgl32.Mat4
    Renderer int
}

type ExplosionRenderer struct {
    MeshRenderers []*GLMeshRenderer
    ParticleCount int
    MaxAge float32
}

func (r *ExplosionRenderer) Render(object *Object) {
    particles, ok := object.RenderData.([]*Particle)
    if !ok {
        panic(fmt.Sprintf("Unexpected RenderData type: %t: %v\n",
            object.RenderData, object.RenderData))
    }

    age := TickTime.TickEnd - object.Creation
    Trace.Printf("Rendering explosion: MaxAge: %v, age: %v, alpha: %v\n", r.MaxAge, age, 1.0 - float32(age) / float32(r.MaxAge))
    for i := 0; i < len(particles); i++ {
        r.MeshRenderers[particles[i].Renderer].RenderAt(
            object.Model.Mul4(particles[i].Model),
            map[int]interface{} {
                UNIFORM_ALPHA: 1.0 - float32(age) / float32(r.MaxAge),
            })
    }
}

func (r *ExplosionRenderer) DebugRender(object *Object) {
    Trace.Printf("ExplosionRenderer: Model matrix:\n%v\n", object.Model)
}

func (r *ExplosionRenderer) Animate(object *Object) {
    particles, ok := object.RenderData.([]*Particle)
    if particles != nil && !ok {
        panic(fmt.Sprintf("Unexpected RenderData type: %t: %v\n",
            object.RenderData, object.RenderData))
    }
    Trace.Printf("Animating particles: %v\n", len(particles))

    age := float32(TickTime.TickEnd - object.Creation)
    if particles == nil {
        particles = make([]*Particle, r.ParticleCount)
        for i:= 0; i < r.ParticleCount; i++ {
            particles[i] = r.createRandomParticle()
            Trace.Printf("Created particle with mesh: %v\n", r.MeshRenderers[particles[i].Renderer].MeshVertices)
        }
    }

    for j := 0; j < len(particles); j++ {
        if age >= r.MaxAge {
            Trace.Printf("Removing particle: before: %v\n", len(particles))
            if len(particles) > 1 {
                particles = append(particles[:j], particles[j+1:]...)
            } else {
                particles = particles[0:0]
            }
            j--
            Trace.Printf("Removing particle: after: %v\n", len(particles))
        } else {
            particles[j].Model = particles[j].Velocity.Mul4(particles[j].Model)
        }
    }

    object.RenderData = particles
}

func (r *ExplosionRenderer) createRandomParticle() *Particle {
    renderer := rand.Intn(len(r.MeshRenderers))
    velocity := mgl32.HomogRotate3DZ(rand.Float32() * 0.01)
    velocity = velocity.Mul4(mgl32.Translate3D(rand.Float32(), rand.Float32(), 0.0))

    return &Particle {
        Renderer: renderer,
        Velocity: velocity,
        Model: mgl32.Ident4(),
    }
}

func InitExplosionRenderer(config *ExplosionRendererConfig) (Renderer, error) {
    meshRenderers := []*GLMeshRenderer {}
    for _, meshRendererConfig := range config.MeshRenderers {
        texture, _, _, err := InitTexture(meshRendererConfig.Texture)
        if err != nil {
            return nil, err
        }

        meshRenderer, err := InitMeshRenderer(
            meshRendererConfig.VertexShader,
            meshRendererConfig.FragmentShader,
            texture,
            []int { UNIFORM_ALPHA },
            meshRendererConfig.MeshVertices)
        if err != nil {
            return nil, err
        }

        meshRenderers = append(meshRenderers, meshRenderer)
    }

    return &ExplosionRenderer {
        ParticleCount: config.ParticleCount,
        MaxAge: config.MaxAge,
        MeshRenderers: meshRenderers,
    }, nil
}
