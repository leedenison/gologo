package gologo

import (
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	// Make sure main thread is locked so that OpenGL calls
	// are always made from the same thread.
	runtime.LockOSThread()
}

func Init() {
	// Use ioutil.Discard to disable
	InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	if err := InitWindow(); err != nil {
		Error.Fatalln("window.Init failed:", err)
	}

	if err := CreateWindow(title); err != nil {
		Error.Fatalln("window.CreateWindow failed:", err)
	}

	UpdateWindowProjection()

	if err := InitRender(); err != nil {
		Error.Fatalln("render.Init failed:", err)
	}

	if err := InitTick(); err != nil {
		Error.Fatalln("physics.Init failed:", err)
	}
}

func Run() {
	// TODO: Add close action for quit key pressed
	for !windowState.Main.ShouldClose() {
		ClearBackBuffer()
		Tick()

		if apiCallback != nil {
			apiCallback(GetTime())
		}

		ResolveContacts(GenerateContacts())
		Render()

		glfw.PollEvents()
	}

	Trace.Println("Exiting.")
	return
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

func SetKeyPressedFunction(f func(int, Key)) {
	keyPressedCallback = f
}

func SetKeyReleasedFunction(f func(int, Key)) {
	keyReleasedCallback = f
}

func IsPressed(key Key) bool {
	state := windowState.Main.GetKey(glfw.Key(key))
	return state == glfw.Press
}

func CreateTaggedContactGenerator(
	tag1 string,
	tag2 string,
	penetration PenetrationResolver,
	post PostContactResolver) {
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
