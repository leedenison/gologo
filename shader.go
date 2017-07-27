package gologo

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
    "strings"
)

type GLShader struct {
    Program uint32
    Projection int32
    Model int32
    Texture int32
    Uniforms map[int]int32
}

func InitShaderProgram(vertexShader string, fragmentShader string) (*GLShader, error) {
    programKey := vertexShader + "," + fragmentShader
    program, programExists := glState.Shaders[programKey]
    if !programExists {
        var err error
        program, err = newProgram(SHADERS[vertexShader], SHADERS[fragmentShader])
        if err != nil {
            return nil, err
        }
        glState.Shaders[programKey] = program
    }

    return program, nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (*GLShader, error) {
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

    return &GLShader {
        Program: program,
        Uniforms: map[int]int32 {},
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
