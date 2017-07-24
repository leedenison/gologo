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
    Shaders map[string]uint32
    Textures map[string]*GLTexture
    NextTextureUnit uint32
    Projection mgl32.Mat4
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

func InitShaderProgram(vertexShader string, textureShader string) (uint32, error) {
    programKey := vertexShader + "," + textureShader
    program, programExists := glState.Shaders[programKey]
    if !programExists {
        var err error
        program, err = newProgram(SHADERS[vertexShader], SHADERS[textureShader])
        if err != nil {
            return 0, err
        }
        glState.Shaders[programKey] = program
    }

    return program, nil
}

func InitTexture(texturePath string) (*GLTexture, uint32, uint32, error) {
    var texture, sizeX, sizeY uint32
    var err error
    result, textureExists := glState.Textures[texturePath]
    if !textureExists {
        texture, sizeX, sizeY, err = newTexture(
            executablePath + PATH_SEPARATOR + texturePath,
            glState.NextTextureUnit)
        if err != nil {
            return nil, 0, 0, err
        }
        result = &GLTexture {
            ID: texture,
            TextureUnit: glState.NextTextureUnit,
        }
        glState.Textures[texturePath] = result
        glState.NextTextureUnit++
    }

    return result, sizeX, sizeY, nil
}
