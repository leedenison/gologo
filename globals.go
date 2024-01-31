package gologo

import (
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
)

/////////////////////////////////////////////////////////////
// API globals
//

var (
	defaultPosition    = mgl32.Vec3{defaultWinSizeX / 2, defaultWinSizeY / 2, 0.0}
	defaultOrientation = 0.0
	defaultScale       = 1.0
)

var apiCallback func(int)

/////////////////////////////////////////////////////////////
// OS globals
//

const (
	resourcePath  = "res"
	pathSeparator = "/"
)

var executablePath string

/////////////////////////////////////////////////////////////
// Window globals
//

const (
	title         = "Gologo!"
	gologoMainWin = "GOLOGO_MAIN"
)

var rendered = []*Object{}

const (
	defaultWinSizeX = 1024
	defaultWinSizeY = 768
)

var (
	keyPressedCallback  func(Key)
	keyReleasedCallback func(Key)
)

var (
	mouseButtonPressedCallback  func(MouseButton)
	mouseButtonReleasedCallback func(MouseButton)
	cursorPositionCallback      func(float64, float64)
)

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

const (
	none   = "NONE"
	circle = "CIRCLE"
)

// TODO(leedenison): We probably don't need SPRITE_CIRCLE after the primitive initialization refactoring
const spriteCircle = "SPRITE_CIRCLE"

const circleMeshSizeFactor = 0.65

const (
	areaToMassRatio      = 0.5
	maxContactIterations = 2
)

var (
	contactGenerators = []ContactGenerator{}
	forceGenerators   = []ForceGenerator{}
)

var integrated = []*Object{}

/////////////////////////////////////////////////////////////
// Object template globals
//

var (
	configs   = map[string]*TemplateConfig{}
	templates = map[string]*Template{}
)

var physicsTypes = map[string]reflect.Type{
	"CIRCLE": reflect.TypeOf(CircleConfig{}),
}

var rendererTypes = map[string]reflect.Type{
	"MESH_RENDERER":      reflect.TypeOf(MeshRendererConfig{}),
	"SPRITE_RENDERER":    reflect.TypeOf(SpriteRendererConfig{}),
	"TEXT_RENDERER":      reflect.TypeOf(TextRendererConfig{}),
	"EXPLOSION_RENDERER": reflect.TypeOf(ExplosionRendererConfig{}),
}
