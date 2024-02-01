// gologo provides a simplified drawing API for learning to code.
//
// gologo uses OpenGL to provide a cross platform drawing API that
// abstracts away the complexities of OpenGL.
//
// A minimal, example gologo program to display a window and wait for it to
// be closed, might be:
//
//	package main
//
//	import (
//		"github.com/leedenison/gologo"
//	)
//
//	func main() {
//		// Create and display the main window
//		g := gologo.Init()
//		defer g.Close()
//
//		// Loop and check if the window has been closed
//		for !g.Window.ShouldClose() {
//			g.CheckForEvents()
//		}
//	}
package gologo

import (
	"os"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/render"
	"github.com/leedenison/gologo/time"
)

type Gologo struct {
	Window *glfw.Window
}

type Config struct {
	Width  int
	Height int
	Title  string
}

const (
	defaultTitle    = "Gologo!"
	defaultWinSizeX = 1024
	defaultWinSizeY = 768
)

func init() {
	// Make sure main thread is locked so that OpenGL calls
	// are always made from the same thread.
	runtime.LockOSThread()
}

func Init() *Gologo {
	return InitWithConfig(
		Config{
			Width:  defaultWinSizeX,
			Height: defaultWinSizeY,
			Title:  defaultTitle,
		})
}

func InitWithConfig(config Config) *Gologo {
	// Use io.Discard to disable
	log.InitLogger(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	if err := glfw.Init(); err != nil {
		log.Error.Fatalln("glfw.Init failed:", err)
	}

	window, err := render.CreateWindow(config.Title, config.Width, config.Height)
	if err != nil {
		log.Error.Fatalln("window.CreateWindow failed:", err)
	}

	width, height := window.GetSize()
	render.Set2DProjection(float32(width), float32(height))

	if err := render.InitOpenGL(); err != nil {
		log.Error.Fatalln("InitOpenGL failed:", err)
	}

	if err := time.InitTick(); err != nil {
		log.Error.Fatalln("InitTick failed:", err)
	}

	return &Gologo{
		Window: window,
	}
}

func (g *Gologo) GetWindowCenter() [2]float32 {
	width, height := g.Window.GetSize()
	return [2]float32{
		float32(width) / 2.0,
		float32(height) / 2.0,
	}
}

func (g *Gologo) ClearBackBuffer() {
	// Clear the OpenGL back buffer
	render.ClearBackBuffer()
}

func (g *Gologo) CheckForEvents() {
	glfw.PollEvents()
}

func (g *Gologo) Close() {
	glfw.Terminate()
}
