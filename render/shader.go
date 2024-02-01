package render

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"

	// Bring in png so we support this file format
	_ "image/png"
)

// GLShader : Stores core info for a GL shader
type GLShader struct {
	Program    uint32
	Projection int32
	Model      int32
	Uniforms   map[int]int32
}

const float32SizeBytes = 4

/////////////////////////////////////////////////////////////
// Shader program globals
//

var (
	UniformProjection = 1
	UniformModel      = 2
	UniformTexture    = 3
	UniformAlpha      = 4
	UniformColor      = 5
)

var (
	shaderUniformLocProjection = gl.Str("projection\x00")
	shaderUniformLocModel      = gl.Str("model\x00")
)

var shaderUniforms = map[int]*uint8{
	UniformTexture: gl.Str("tex\x00"),
	UniformAlpha:   gl.Str("alpha\x00"),
	UniformColor:   gl.Str("color\x00"),
}

var fragLocOutputColor = gl.Str("outputColor\x00")

var (
	attribLocVertex         = gl.Str("vert\x00")
	attribLocVertexTexCoord = gl.Str("vertTexCoord\x00")
)

var shaders = map[string]string{
	"FULLSCREEN_VERTEX_SHADER": `
#version 330

in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(vert, 1);
}
` + "\x00",
	"ORTHO_VERTEX_SHADER": `
#version 330

uniform mat4 projection;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * model * vec4(vert, 1);
}
` + "\x00",

	"TEXTURE_FRAGMENT_SHADER": `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;
out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00",

	"ALPHA_FRAGMENT_SHADER": `
#version 330

uniform sampler2D tex;
uniform float alpha;

vec4 texColor;
in vec2 fragTexCoord;
out vec4 outputColor;

void main() {
    texColor = texture(tex, fragTexCoord);
    outputColor = vec4(texColor.rgb, texColor.a * alpha);
}
` + "\x00",

	"COLOR_FRAGMENT_SHADER": `
#version 330

uniform vec4 color;

out vec4 outputColor;

void main() {
    outputColor = color;
}
` + "\x00",
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
