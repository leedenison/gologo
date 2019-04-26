package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

const float32SizeBytes = 4

/////////////////////////////////////////////////////////////
// Rendering globals
//

const GlMeshStride = 5
const glMeshStrideBytes = GlMeshStride * float32SizeBytes

var CreateMeshRenderer func(
	vertexShader string,
	fragmentShader string,
	uniforms []int,
	uniformValues map[int]interface{},
	meshVertices []float32) (*MeshRenderer, error)

var CreateTexture func(texturePath string) (*GLTexture, error)

var glState = GLState{
	Shaders:         map[string]*GLShader{},
	Textures:        map[string]*GLTexture{},
	NextTextureUnit: gl.TEXTURE0,
}

/////////////////////////////////////////////////////////////
// Shader program globals
//

var UniformProjection = 1
var UniformModel = 2
var UniformTexture = 3
var UniformAlpha = 4
var UniformColor = 5

var shaderUniformLocProjection = gl.Str("projection\x00")
var shaderUniformLocModel = gl.Str("model\x00")

var shaderUniforms = map[int]*uint8{
	UniformTexture: gl.Str("tex\x00"),
	UniformAlpha:   gl.Str("alpha\x00"),
	UniformColor:   gl.Str("color\x00"),
}

var fragLocOutputColor = gl.Str("outputColor\x00")

var attribLocVertex = gl.Str("vert\x00")
var attribLocVertexTexCoord = gl.Str("vertTexCoord\x00")

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
