package gologo

import (
    "os"
    "path/filepath"
    "runtime"
    "sort"

    "github.com/go-gl/gl/v4.1-core/gl"
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/mathgl/mgl32"
)

func init() {
    // Use ioutil.Discard to disable
    InitLogging(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

    // Make sure main thread is locked so that OpenGL calls
    // are always made from the same thread.
    runtime.LockOSThread()

    path, err := getProgramPath()
    if err != nil {
        Error.Fatalln("Failed to determine executable path:", err)
    }

    executablePath = path

    // TODO: Set up model matrices and rigid body physics data for defined objects
    ObjectTypeConfigs, err = LoadObjectTypeConfigs(
        executablePath + PATH_SEPARATOR + RESOURCE_PATH)
    if err != nil {
        Error.Fatalln("Failed to load resources:", err)
    }

    TickTime.Start = 0
    TickTime.TickEnd = 0
}

func getProgramPath() (string, error) {
    ex, err := os.Executable()
    if err != nil {
        return "", err
    }

    return filepath.Dir(ex), nil
}

func Run(title string) {
    glWin.Width = DEFAULT_WIN_SIZE_X
    glWin.Height = DEFAULT_WIN_SIZE_Y

    if err := glfw.Init(); err != nil {
        Error.Fatalln("glfw.Init failed:", err)
    }
    defer glfw.Terminate()

    glfw.WindowHint(glfw.Resizable, glfw.False)
    glfw.WindowHint(glfw.ContextVersionMajor, 4)
    glfw.WindowHint(glfw.ContextVersionMinor, 1)
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
    glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

    window, err := glfw.CreateWindow(glWin.Width, glWin.Height, title, nil, nil)
    if err != nil {
        Error.Fatalln("glfw.CreateWindow failed:", err)
    }

    glWin.Window = window
    glWin.Window.MakeContextCurrent()

    // Initialize Glow
    if err := gl.Init(); err != nil {
        Error.Fatalln("gl.Init failed:", err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    Trace.Println("OpenGL version:", version)

    // setup logical world size and use mgl32.Ortho2D to set up bounding
    // volume, logical world size should be the same aspect ratio as the window
    glState.Projection = mgl32.Ortho2D(
        0, float32(glWin.Width), 0, float32(glWin.Height))

    // Set up shaders for defined object types
    ObjectTypes, err = InitObjectTypes(ObjectTypeConfigs)
    if err != nil {
        Error.Fatalln("Failed to initialize object types:", err)
    }

    // Update objects with their intialised object type data
    err = InitObjects(Objects, ObjectTypes)
    if err != nil {
        Error.Fatalln("Failed to initialize objects:", err)
    }

    glWin.Window.SetKeyCallback(KeyCallback)

    // Configure global settings
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)

    TickTime.Start = glfw.GetTime()
    TickTime.TickEnd = TickTime.Start

    // TODO: Add close action for quit key pressed
    for !window.ShouldClose() {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        time := glfw.GetTime()
        TickTime.TickInterval = time - TickTime.TickEnd
        TickTime.TickEnd = time

        if apiCallback != nil {
            apiCallback(PhysicsTime())
        }

        PhysicsTick()

        // Sort the objects by the rendering zorder
        // TODO: Add a dirty flag to know if objects has been modified
        sort.Sort(ByZOrder(Objects))
        for _, object := range Objects {
            object.ObjectType.Renderer.Animate(object)
            //object.ObjectType.Renderer.DebugRender(object)
            object.ObjectType.Renderer.Render(object)
        }

        window.SwapBuffers()
        glfw.PollEvents()
    }

    Info.Println("Exiting.")
    return
}
