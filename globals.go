package gologo

import (
    "reflect"
    "github.com/go-gl/mathgl/mgl32"
    "github.com/go-gl/gl/v4.1-core/gl"
)

/////////////////////////////////////////////////////////////
// API globals
//

var DEFAULT_POSITION = mgl32.Translate3D(
    DEFAULT_WIN_SIZE_X / 2,
    DEFAULT_WIN_SIZE_Y / 2,
    0.0)

var DEFAULT_ORIENTATION = mgl32.Ident4()
var DEFAULT_SCALE = mgl32.Ident4()

var apiCallback func(int)

/////////////////////////////////////////////////////////////
// OS globals
//

const FLOAT32_SIZE_BYTES = 4

const RESOURCE_PATH = "res"
const PATH_SEPARATOR = "/"

var executablePath string

/////////////////////////////////////////////////////////////
// Window globals
//

const TITLE = "Gologo!"
const GOLOGO_MAIN_WIN = "GOLOGO_MAIN"

const DEFAULT_WIN_SIZE_X = 1024
const DEFAULT_WIN_SIZE_Y = 768

var windowState = WindowState {
    Width: 0,
    Height: 0,
}

var keyPressedCallback func(int, Key)
var keyReleasedCallback func(int, Key)

/////////////////////////////////////////////////////////////
// Lifecycle globals
//

var process = &Lifecycle {}

/////////////////////////////////////////////////////////////
// Tags globals
//

var tags = map[string]ObjectSet {}

/////////////////////////////////////////////////////////////
// Physics globals
//

const NONE = "NONE"
const CIRCLE = "CIRCLE"
const SPRITE_CIRCLE = "SPRITE_CIRCLE"

const CIRCLE_MESH_SIZE_FACTOR = 0.65

const AREA_TO_MASS_RATIO = 0.5
const MAX_CONTACT_ITERATIONS = 2

var tick = TickState {}

var contactGenerators = []ContactGenerator {}

/////////////////////////////////////////////////////////////
// Object template globals
//

var configs = map[string]*TemplateConfig {}
var templates = map[string]*Template {}

var rendererTypes = map[string]reflect.Type {
    "MESH_RENDERER": reflect.TypeOf(MeshRendererConfig {}),
    "SPRITE_RENDERER": reflect.TypeOf(SpriteRendererConfig {}),
    "TEXT_RENDERER": reflect.TypeOf(TextRendererConfig {}),
    "EXPLOSION_RENDERER": reflect.TypeOf(ExplosionRendererConfig {}),
}

var physicsTypes = map[string]reflect.Type {
    "CIRCLE": reflect.TypeOf(CircleConfig {}),
}

/////////////////////////////////////////////////////////////
// Rendering globals
//

const GL_MESH_STRIDE = 5
const GL_MESH_STRIDE_BYTES = GL_MESH_STRIDE * FLOAT32_SIZE_BYTES

var glState = GLState {
    Shaders: map[string]*GLShader {},
    Textures: map[string]*GLTexture {},
    NextTextureUnit: gl.TEXTURE0,
}

var rendered = []*Object {}

/////////////////////////////////////////////////////////////
// Shader program globals
//

var UNIFORM_PROJECTION = 1
var UNIFORM_MODEL = 2
var UNIFORM_TEXTURE = 3
var UNIFORM_ALPHA = 4
var UNIFORM_COLOR = 5

var UNIFORM_LOC_PROJECTION = gl.Str("projection\x00")
var UNIFORM_LOC_MODEL = gl.Str("model\x00")

var UNIFORMS = map[int]*uint8 {
    UNIFORM_TEXTURE: gl.Str("tex\x00"),
    UNIFORM_ALPHA: gl.Str("alpha\x00"),
    UNIFORM_COLOR: gl.Str("color\x00"),
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

    "COLOR_FRAGMENT_SHADER": `
#version 330

uniform vec4 color;

out vec4 outputColor;

void main() {
    outputColor = color;
}
` + "\x00",
}
