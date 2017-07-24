package gologo

import (
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
    Shader uint32
    ProjectionUniform int32
    ModelUniform int32
    TextureUniform int32
    Texture *GLTexture
    Mesh uint32
    MeshVertices []float32
    VertexCount int32
}

func (r *GLMeshRenderer) Render(object *Object) {
    r.RenderAt(object.Model)
}

func (r *GLMeshRenderer) RenderAt(model mgl32.Mat4) {
    gl.UseProgram(r.Shader)
    gl.UniformMatrix4fv(r.ModelUniform, 1, false, &model[0])

    gl.BindVertexArray(r.Mesh)

    gl.ActiveTexture(gl.TEXTURE0)
    gl.Uniform1i(r.TextureUniform, 0)
    gl.BindTexture(gl.TEXTURE_2D, r.Texture.ID)

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
        config.TextureShader,
        texture,
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
        config.TextureShader,
        texture,
        config.MeshVertices)
}

func InitMeshRenderer(
        vertexShader string,
        textureShader string,
        texture *GLTexture,
        meshVertices []float32) (Renderer, error) {
    err := ValidateMeshRenderConfig(vertexShader, textureShader, meshVertices)
    if err != nil {
        return nil, err
    }

    program, err := InitShaderProgram(vertexShader, textureShader)
    if err != nil {
        return nil, err
    }

    gl.UseProgram(program)

    projectionU, modelU, textureU := InitUniforms(program)
    gl.UniformMatrix4fv(projectionU, 1, false, &glState.Projection[0])

    mesh := InitObjectMesh(program, meshVertices)

    return &GLMeshRenderer {
        Shader: program,
        ProjectionUniform: projectionU,
        ModelUniform: modelU,
        TextureUniform: textureU,
        Texture: texture,
        Mesh: mesh,
        MeshVertices: meshVertices,
        VertexCount: int32(len(meshVertices) / GL_MESH_STRIDE),
    }, nil
}

func InitUniforms(shader uint32) (int32, int32, int32) {
    // Bind the uniform variables
    projectionUniform := gl.GetUniformLocation(shader, gl.Str("projection\x00"))
    modelUniform := gl.GetUniformLocation(shader, gl.Str("model\x00"))
    textureUniform := gl.GetUniformLocation(shader, gl.Str("tex\x00"))

    gl.BindFragDataLocation(shader, 0, gl.Str("outputColor\x00"))

    return projectionUniform, modelUniform, textureUniform
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

    vertAttrib := uint32(gl.GetAttribLocation(shader, gl.Str("vert\x00")))
    gl.EnableVertexAttribArray(vertAttrib)
    gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, GL_MESH_STRIDE_BYTES,
        gl.PtrOffset(0))

    texCoordAttrib := uint32(gl.GetAttribLocation(shader, gl.Str("vertTexCoord\x00")))
    gl.EnableVertexAttribArray(texCoordAttrib)
    gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, GL_MESH_STRIDE_BYTES,
        gl.PtrOffset(3*4))

    return vao
}

func ValidateMeshRenderConfig(
        vertexShader string,
        textureShader string,
        meshVertices []float32) error {
    if vertexShader == "" {
        return errors.New("Missing required field: 'VertexShader'")
    } else if _, ok := SHADERS[vertexShader]; !ok {
        return errors.Errorf("Unknown 'VertexShader': %v", vertexShader)
    }

    if textureShader == "" {
        return errors.New("Missing required field: 'TextureShader'")
    } else if _, ok := SHADERS[textureShader]; !ok {
        return errors.Errorf("Unknown 'TextureShader': %v", textureShader)
    }

    if len(meshVertices) == 0 {
        return errors.New("Missing required field: 'MeshVertices'")
    }

    return nil
}

