package gologo

import (
	"os"
	"path/filepath"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	windowState.Width = defaultWinSizeX
	windowState.Height = defaultWinSizeY

	_, err := GetExecutablePath()
	if err != nil {
		Error.Fatalln("Failed to determine executable path:", err)
	}
}

type WindowState struct {
	Main   *glfw.Window
	Width  int
	Height int
}

func InitWindow() error {
	if err := glfw.Init(); err != nil {
		Error.Println("glfw.Init failed:", err)
		return err
	}

	process.RegisterCleanup(Cleanup)
	return nil
}

func Cleanup() {
	glfw.Terminate()
}

func CreateWindow(title string) error {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowState.Width, windowState.Height, title, nil, nil)
	if err != nil {
		Error.Println("glfw.CreateWindow failed:", err)
	}

	windowState.Main = window
	windowState.Main.MakeContextCurrent()
	windowState.Main.SetKeyCallback(KeyCallback)

	return nil
}

func GetWindowSize() [2]float32 {
	return [2]float32{float32(windowState.Height), float32(windowState.Width)}
}

func GetWindowCenter() [2]float32 {
	return [2]float32{
		float32(windowState.Width) / 2.0,
		float32(windowState.Height) / 2.0,
	}
}

func GetResourcePath() (string, error) {
	path, err := GetExecutablePath()
	if err != nil {
		return "", err
	}

	return path + pathSeparator + resourcePath, nil
}

func GetExecutablePath() (string, error) {
	if executablePath == "" {
		ex, err := os.Executable()
		if err != nil {
			return "", err
		}
		executablePath = filepath.Dir(ex)
	}

	return executablePath, nil
}

func KeyCallback(
	window *glfw.Window,
	key glfw.Key,
	scancode int,
	action glfw.Action,
	mods glfw.ModifierKey) {
	if action == glfw.Press && key == glfw.KeyEscape {
		window.SetShouldClose(true)
	} else if keyPressedCallback != nil && action == glfw.Press {
		keyPressedCallback(GetTickTime(), Key(key))
	} else if keyReleasedCallback != nil && action == glfw.Release {
		keyReleasedCallback(GetTickTime(), Key(key))
	}
}
