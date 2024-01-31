package gologo

import (
	"os"
	"runtime"
	"sort"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/opengl"
	"github.com/leedenison/gologo/time"
)

func init() {
	// Make sure main thread is locked so that OpenGL calls
	// are always made from the same thread.
	runtime.LockOSThread()
}

func Init() {
	// Use ioutil.Discard to disable
	log.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	if err := InitWindow(); err != nil {
		log.Error.Fatalln("window.Init failed:", err)
	}

	if err := CreateWindow(title); err != nil {
		log.Error.Fatalln("window.CreateWindow failed:", err)
	}

	opengl.UpdateProjection(float32(windowState.Width), float32(windowState.Height))

	if err := opengl.InitOpenGL(); err != nil {
		log.Error.Fatalln("opengl.InitOpenGL failed:", err)
	}

	if err := time.InitTick(); err != nil {
		log.Error.Fatalln("time.InitTick failed:", err)
	}
}

func Run() {
	for !ShouldCloseWindow() {
		ClearWindow()
		time.Tick()

		if apiCallback != nil {
			apiCallback(time.GetTime())
		}

		ClearForces(integrated)
		GenerateForces(time.TimeState.Interval)
		Integrate(time.TimeState.Interval)
		ResolveContacts(GenerateContacts())
		Draw()

		CheckForUserInput()
	}

	log.Trace.Println("Exiting.")
}

func ShouldCloseWindow() bool {
	return windowState.Main.ShouldClose()
}

func ClearWindow() {
	opengl.ClearBackBuffer()
}

func CheckForUserInput() {
	glfw.PollEvents()
}

func GetScreenWidth() int {
	return int(windowState.Width)
}

func GetScreenHeight() int {
	return int(windowState.Height)
}

func SetTickFunction(f func(int)) {
	apiCallback = f
}

func SetKeyPressedFunction(f func(Key)) {
	keyPressedCallback = f
}

func SetKeyReleasedFunction(f func(Key)) {
	keyReleasedCallback = f
}

func IsPressed(key Key) bool {
	state := windowState.Main.GetKey(glfw.Key(key))
	return state == glfw.Press
}

func SetMouseButtonPressedFunction(f func(MouseButton)) {
	mouseButtonPressedCallback = f
}

func SetMouseButtonReleasedFunction(f func(MouseButton)) {
	mouseButtonReleasedCallback = f
}

func SetCursorPositionFunction(f func(float64, float64)) {
	cursorPositionCallback = f
}

func GetCursorPosition() (float64, float64) {
	return windowState.Main.GetCursorPos()
}

func CreateTaggedContactGenerator(
	tag1 string,
	tag2 string,
	penetration PenetrationResolver,
	post PostContactResolver,
) {
	contactGenerators = append(contactGenerators,
		&TaggedContactGenerator{
			SourceTag:           tag1,
			TargetTag:           tag2,
			PenetrationResolver: penetration,
			PostContactResolver: post,
		})
}

func SetScreenEdgeTag(tag string) {
	contactGenerators = append(contactGenerators,
		&ScreenEdgeContactGenerator{
			Tag:                 tag,
			PenetrationResolver: &PerpendicularPenetrationResolver{},
		})
}

// TODO(leedenison): Add duration to render for animate.
func Draw() {
	// Sort the objects by the rendering zorder
	sort.Sort(ByZOrder(rendered))
	for _, object := range rendered {
		object.Renderer.Animate(object.GetModel())
		// object.Renderer.DebugRender(object.GetModel())
		object.Renderer.Render(object.GetModel())
	}
	windowState.Main.SwapBuffers()
}

func TagRender(object *Object) {
	rendered = append(rendered, object)
}

func UntagRender(object *Object) {
	for i := 0; i < len(rendered); i++ {
		if object == rendered[i] {
			if len(rendered) > 1 {
				rendered = append(rendered[:i], rendered[i+1:]...)
			} else {
				rendered = rendered[0:0]
			}
			i--
		}
	}
}
