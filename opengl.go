package gologo

import (
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/mathgl/mgl32"
)

type WindowState struct {
    Window *glfw.Window
    Width int
    Height int
}

type GLState struct {
    Shaders map[string]*GLShader
    Textures map[string]*GLTexture
    NextTextureUnit uint32
    Projection mgl32.Mat4
}

type GLTexture struct {
    ID uint32
    TextureUnit uint32
    Size [2]uint32
}

type TimeState struct {
    Start float64
    TickEnd float64
    TickInterval float64
}

func KeyCallback(
        window *glfw.Window,
        key glfw.Key,
        scancode int,
        action glfw.Action,
        mods glfw.ModifierKey) {
    if keyPressedCallback != nil && action == glfw.Press {
        keyPressedCallback(PhysicsTime(), Key(key))
    } else if keyReleasedCallback != nil && action == glfw.Release {
		keyReleasedCallback(PhysicsTime(), Key(key))
	}
}

func PhysicsTime() int {
    return int(1000 * (glfw.GetTime() - TickTime.Start))
}

func InitTexture(texturePath string) (*GLTexture, error) {
    result, textureExists := glState.Textures[texturePath]
    if !textureExists {
        texture, sizeX, sizeY, err := newTexture(
            executablePath + PATH_SEPARATOR + texturePath,
            glState.NextTextureUnit)
        if err != nil {
            return nil, err
        }
        result = &GLTexture {
            ID: texture,
            TextureUnit: glState.NextTextureUnit,
            Size: [2]uint32 { sizeX, sizeY },
        }
        glState.Textures[texturePath] = result
        glState.NextTextureUnit++
    }

    return result, nil
}
