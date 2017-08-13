package gologo

import (
    "os"
    "fmt"
    "github.com/pkg/errors"
    "github.com/go-gl/gl/v4.1-core/gl"
    "io/ioutil"
    "image"
    "image/draw"
    _ "image/png"
    "encoding/json"
    "regexp"
)


/////////////////////////////////////////////////////////////
// Config types
//

type ObjectTypeConfig struct {
    Name string
    RendererType string
    Renderer json.RawMessage
    RendererConfig interface{}
    PhysicsPrimitiveType string
    PhysicsPrimitive json.RawMessage
    PhysicsPrimitiveConfig interface{}
}

type GLMeshRendererConfig struct {
    VertexShader string
    FragmentShader string
    Texture string
    MeshVertices []float32
}

type SpriteMeshRendererConfig struct {
    VertexShader string
    FragmentShader string
    Texture string
    TextureOrigin []int32
    MeshScaling float32
}

type ExplosionRendererConfig struct {
    ParticleCount int
    MaxAge float32
    MeshRenderers []GLMeshRendererConfig
}

type TextRendererConfig struct {
    MeshRenderers map[string]CharRendererConfig
    CharSpacer float32
}

type CharRendererConfig struct {
    VertexShader string
    FragmentShader string
    Texture string
    TextureSize [2]float32
    TextureRect [][2]float32
    CharRect [][2]float32
}

type CircleConfig struct {
    Radius float32
    InverseMass float32
}

type SpriteCirclePrimitiveConfig struct {}

// LoadObjectTypes loads Gologo object type data from a resource directory.
func LoadObjectTypeConfigs(resourceDir string) (map[string]*ObjectTypeConfig, error) {
    result := map[string]*ObjectTypeConfig {}

    files, err := ioutil.ReadDir(resourceDir)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to load resources.")
    }

    for _, file := range files {
        Trace.Printf("Config file: %v\n", file.Name())
        matched, _ := regexp.MatchString(".*\\.json$", file.Name())
        if file.IsDir() || !matched {
            continue
        }

        filePath := resourceDir + "/" + file.Name()
        objectTypeConfig, err := LoadObjectTypeConfig(filePath)
        if err != nil {
            Warning.Println("Skipping resource:", err)
            continue
        }

        if objectTypeConfig.Name == "" {
            return nil, errors.New("ObjectType is missing required field: 'Name'")
        }

        result[objectTypeConfig.Name] = objectTypeConfig
    }

    return result, nil
}

// LoadObjectType loads a Gologo object type from the specified path.
func LoadObjectTypeConfig(resourcePath string) (*ObjectTypeConfig, error) {
    parseResult := ObjectTypeConfig {}

    resourceJson, err := ioutil.ReadFile(resourcePath)
    if err != nil {
        return nil, errors.Wrapf(err, "Failed to load resource: %s", resourcePath)
    }

    err = json.Unmarshal(resourceJson, &parseResult)
    if err != nil {
        return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
    }

    switch parseResult.RendererType {
    case GL_MESH_RENDERER:
        rendererConfig := GLMeshRendererConfig {}
        err = json.Unmarshal(parseResult.Renderer, &rendererConfig)
        if err != nil {
            return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
        }
        parseResult.RendererConfig = rendererConfig
    case SPRITE_MESH_RENDERER:
        rendererConfig := SpriteMeshRendererConfig {}
        err = json.Unmarshal(parseResult.Renderer, &rendererConfig)
        if err != nil {
            return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
        }
        parseResult.RendererConfig = rendererConfig
    case EXPLOSION_RENDERER:
        rendererConfig := ExplosionRendererConfig {}
        err = json.Unmarshal(parseResult.Renderer, &rendererConfig)
        if err != nil {
            return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
        }
        parseResult.RendererConfig = rendererConfig
    case TEXT_RENDERER:
        rendererConfig := TextRendererConfig {}
        err = json.Unmarshal(parseResult.Renderer, &rendererConfig)
        if err != nil {
            return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
        }
        parseResult.RendererConfig = rendererConfig
    default:
        return nil, errors.Errorf("Unknown RenderType: %v\n", parseResult.RendererType)
    }

    switch parseResult.PhysicsPrimitiveType {
    case NONE:
        break
    case SPRITE_CIRCLE:
        break
    case CIRCLE:
        physicsPrimitiveConfig := CircleConfig {}
        err = json.Unmarshal(parseResult.Renderer, &physicsPrimitiveConfig)
        if err != nil {
            return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
        }
        parseResult.PhysicsPrimitiveConfig = physicsPrimitiveConfig
    default:
        return nil, errors.Errorf("Unknown PhysicsPrimitiveType: %v\n",
            parseResult.PhysicsPrimitiveType)
    }

    return &parseResult, nil
}

func newTexture(file string, textureUnit uint32) (uint32, uint32, uint32, error) {
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
