package gologo

import (
    "github.com/go-gl/gl/v4.1-core/gl"
)

/////////////////////////////////////////////////////////////
// API globals
//

var apiCallback func(int)

var keyPressedCallback func(int, Key)
var keyReleasedCallback func(int, Key)

/////////////////////////////////////////////////////////////
// OS globals
//

const PATH_SEPARATOR = "/"
const RESOURCE_PATH = "res"

const FLOAT32_SIZE_BYTES = 4

var executablePath string

/////////////////////////////////////////////////////////////
// Physics globals
//

const CIRCLE = "CIRCLE"
const SPRITE_CIRCLE = "SPRITE_CIRCLE"
const CIRCLE_MESH_SIZE_FACTOR = 0.85

const AREA_TO_MASS_RATIO = 0.5
const MAX_CONTACT_ITERATIONS = 2

var ObjectTypeConfigs = map[string]*ObjectTypeConfig {}
var ObjectTypes = map[string]*ObjectType {}
var Objects = []*Object {}
var Tags = map[string]ObjectSet {}

var contactGenerators = []ContactGenerator {}

/////////////////////////////////////////////////////////////
// Rendering globals
//

const GOLOGO_MAIN_WIN = "GOLOGO_MAIN"

const GL_MESH_RENDERER = "GL_MESH_RENDERER"
const SPRITE_MESH_RENDERER = "SPRITE_MESH_RENDERER"
const EXPLOSION_RENDERER = "EXPLOSION_RENDERER"

const DEFAULT_WIN_SIZE_X = 1024
const DEFAULT_WIN_SIZE_Y = 768

const GL_MESH_STRIDE = 5
const GL_MESH_STRIDE_BYTES = GL_MESH_STRIDE * FLOAT32_SIZE_BYTES

var glWin = WindowState {
    Width: 0,
    Height: 0,
}

var glState = GLState {
    Shaders: map[string]*GLShader {},
    Textures: map[string]*GLTexture {},
    NextTextureUnit: gl.TEXTURE0,
}

var TickTime = TimeState {}

/////////////////////////////////////////////////////////////
// Shader program globals
//

var UNIFORM_PROJECTION = 1
var UNIFORM_MODEL = 2
var UNIFORM_TEXTURE = 3
var UNIFORM_ALPHA = 4

var UNIFORM_LOC_PROJECTION = gl.Str("projection\x00")
var UNIFORM_LOC_MODEL = gl.Str("model\x00")
var UNIFORM_LOC_TEXTURE = gl.Str("tex\x00")

var UNIFORMS = map[int]*uint8 {
    UNIFORM_ALPHA: gl.Str("alpha\x00"),
}

var FRAG_LOC_OUTPUT_COLOR = gl.Str("outputColor\x00")

var ATTRIB_LOC_VERTEX = gl.Str("vert\x00")
var ATTRIB_LOC_VERTEX_TEX_COORD = gl.Str("vertTexCoord\x00")

var SHADERS = map[string]string {
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
}
