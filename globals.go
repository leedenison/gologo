package gologo

import (
	"reflect"

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

const resourcePath = "res"
const pathSeparator = "/"

var executablePath string

/////////////////////////////////////////////////////////////
// Window globals
//

const title = "Gologo!"
const gologoMainWin = "GOLOGO_MAIN"

var rendered = []*Object{}

const defaultWinSizeX = 1024
const defaultWinSizeY = 768

var keyPressedCallback func(int, Key)
var keyReleasedCallback func(int, Key)

var windowState = WindowState{
	Width:  0,
	Height: 0,
}

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

// TODO(leedenison): We probably don't need SPRITE_CIRCLE after the primitive initialization refactoring
const spriteCircle = "SPRITE_CIRCLE"

const circleMeshSizeFactor = 0.65

const areaToMassRatio = 0.5
const maxContactIterations = 2

var contactGenerators = []ContactGenerator{}
var forceGenerators = []ForceGenerator{}

var integrated = []*Object{}

/////////////////////////////////////////////////////////////
// Object template globals
//

var configs = map[string]*TemplateConfig{}
var templates = map[string]*Template{}

var physicsTypes = map[string]reflect.Type{
	"CIRCLE": reflect.TypeOf(CircleConfig{}),
}

var rendererTypes = map[string]reflect.Type{
	"MESH_RENDERER":      reflect.TypeOf(MeshRendererConfig{}),
	"SPRITE_RENDERER":    reflect.TypeOf(SpriteRendererConfig{}),
	"TEXT_RENDERER":      reflect.TypeOf(TextRendererConfig{}),
	"EXPLOSION_RENDERER": reflect.TypeOf(ExplosionRendererConfig{}),
}
