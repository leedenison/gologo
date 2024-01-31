package gologo

import (
	"os"
	"path/filepath"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/leedenison/gologo/log"
)

const (
	resourcePath  = "res"
	pathSeparator = "/"
)

var executablePath string

const (
	mainWindow      = "GOLOGO_MAIN"
	defaultTitle    = "Gologo!"
	defaultWinSizeX = 1024
	defaultWinSizeY = 768
)

func init() {
	_, err := GetExecutablePath()
	if err != nil {
		log.Error.Fatalln("Failed to determine executable path:", err)
	}
}

func CreateWindow(title string, width int, height int) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()
	window.SetKeyCallback(KeyCallback)

	return window, nil
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
	mod glfw.ModifierKey,
) {
	if action == glfw.Press && key == glfw.KeyEscape {
		window.SetShouldClose(true)
	}
}
