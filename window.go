package gologo

import (
	"os"
	"path/filepath"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/time"
)

type WindowState struct {
	Main   *glfw.Window
	Width  int
	Height int
}

func init() {
	windowState.Width = defaultWinSizeX
	windowState.Height = defaultWinSizeY

	_, err := GetExecutablePath()
	if err != nil {
		log.Error.Fatalln("Failed to determine executable path:", err)
	}
}

func InitWindow() error {
	if err := glfw.Init(); err != nil {
		log.Error.Println("glfw.Init failed:", err)
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
		log.Error.Println("glfw.CreateWindow failed:", err)
	}

	windowState.Main = window
	windowState.Main.MakeContextCurrent()
	windowState.Main.SetKeyCallback(KeyCallback)

	return nil
}

func GetWindowSize() [2]float32 {
	return [2]float32{float32(windowState.Width), float32(windowState.Height)}
}

func SetWindowSize(s [2]int) {
	windowState.Width = s[0]
	windowState.Height = s[1]
}

func GetWindowCenter() [2]float32 {
	return [2]float32{
		float32(windowState.Width) / 2.0,
		float32(windowState.Height) / 2.0,
	}
}

func GetResourcePath() string {
	return resourcePath
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
		keyPressedCallback(time.GetTickTime(), Key(key))
	} else if keyReleasedCallback != nil && action == glfw.Release {
		keyReleasedCallback(time.GetTickTime(), Key(key))
	}
}
