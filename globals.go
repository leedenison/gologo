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

const DEFAULT_WIN_SIZE_X = 1024
const DEFAULT_WIN_SIZE_Y = 768

const GL_MESH_STRIDE = 5
const GL_MESH_STRIDE_BYTES = GL_MESH_STRIDE * FLOAT32_SIZE_BYTES

var glWin = WindowState {
    Width: 0,
    Height: 0,
}

var glState = GLState {
    Shaders: map[string]uint32 {},
    Textures: map[string]*GLTexture {},
    NextTextureUnit: gl.TEXTURE0,
}

var TickTime = TimeState {}

/////////////////////////////////////////////////////////////
// Shader program globals
//

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

out vec4 outputColor;

void main() {
    outputColor = vec4(1.0, 0.0, 0.0, 1.0);
}
` + "\x00",
}
