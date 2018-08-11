package gologo

import (
	"reflect"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

/////////////////////////////////////////////////////////////
// API globals
//

var defaultPosition = mgl32.Vec3{defaultWinSizeX / 2, defaultWinSizeY / 2, 0.0}
var defaultOrientation = 0.0
var defaultScale = 1.0

var apiCallback func(int)

/////////////////////////////////////////////////////////////
// OS globals
//

const float32SizeBytes = 4

const resourcePath = "res"
const pathSeparator = "/"

var executablePath string

/////////////////////////////////////////////////////////////
// Window globals
//

const title = "Gologo!"
const gologoMainWin = "GOLOGO_MAIN"

const defaultWinSizeX = 1024
const defaultWinSizeY = 768

var windowState = WindowState{
	Width:  0,
	Height: 0,
}

var keyPressedCallback func(int, Key)
var keyReleasedCallback func(int, Key)

/////////////////////////////////////////////////////////////
// Lifecycle globals
//

var process = &Lifecycle{}

/////////////////////////////////////////////////////////////
// Tags globals
//

var tags = map[string]ObjectSet{}

/////////////////////////////////////////////////////////////
// Physics globals
//

const none = "NONE"
const circle = "CIRCLE"
const spriteCircle = "SPRITE_CIRCLE"

const circleMeshSizeFactor = 0.65

const areaToMassRatio = 0.5
const maxContactIterations = 2

var tick = TickState{}

var contactGenerators = []ContactGenerator{}

/////////////////////////////////////////////////////////////
// Object template globals
//

var configs = map[string]*TemplateConfig{}
var templates = map[string]*Template{}

var rendererTypes = map[string]reflect.Type{
	"MESH_RENDERER":      reflect.TypeOf(MeshRendererConfig{}),
	"SPRITE_RENDERER":    reflect.TypeOf(SpriteRendererConfig{}),
	"TEXT_RENDERER":      reflect.TypeOf(TextRendererConfig{}),
	"EXPLOSION_RENDERER": reflect.TypeOf(ExplosionRendererConfig{}),
}

var physicsTypes = map[string]reflect.Type{
	"CIRCLE": reflect.TypeOf(CircleConfig{}),
}

/////////////////////////////////////////////////////////////
// Rendering globals
//

const glMeshStride = 5
const glMeshStrideBytes = glMeshStride * float32SizeBytes

var glState = GLState{
	Shaders:         map[string]*GLShader{},
	Textures:        map[string]*GLTexture{},
	NextTextureUnit: gl.TEXTURE0,
}

var rendered = []*Object{}

/////////////////////////////////////////////////////////////
// Shader program globals
//

var uniformProjection = 1
var uniformModel = 2
var uniformTexture = 3
var uniformAlpha = 4
var uniformColor = 5

var shaderUniformLocProjection = gl.Str("projection\x00")
var shaderUniformLocModel = gl.Str("model\x00")

var shaderUniforms = map[int]*uint8{
	uniformTexture: gl.Str("tex\x00"),
	uniformAlpha:   gl.Str("alpha\x00"),
	uniformColor:   gl.Str("color\x00"),
}

var fragLocOutputColor = gl.Str("outputColor\x00")

var attribLocVertex = gl.Str("vert\x00")
var attribLocVertexTexCoord = gl.Str("vertTexCoord\x00")

var shaders = map[string]string{
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
