package gologowin

import "github.com/AllenDang/w32"
import "github.com/leedenison/gologo/w32ext"
import "unsafe"

func CreateWindowClass(hInstance w32.HINSTANCE, lpszClassName *uint16, wndProc uintptr) w32.WNDCLASSEX {
	var wcex w32.WNDCLASSEX

	// Size of the window object.
	wcex.Size = uint32(unsafe.Sizeof(wcex))

	wcex.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	// Application loop handler procedure.
	wcex.WndProc = wndProc
	// Additional bytes to allocate for the window class struct.
	wcex.ClsExtra = 0
	// Additional bytes to allocation for the window instance struct.
	// If an application uses WNDCLASS to register a dialog box created
	// by using the CLASS directive in the resource file, it must set this
	// member to DLGWINDOWEXTRA.
	wcex.WndExtra = 0
	// A handle to the instance that contains the window procedure for the class.
	wcex.Instance = hInstance

	// Use default IDI_APPLICATION icon.
	wcex.Icon = w32.LoadIcon(hInstance, w32ext.MakeIntResource(w32.IDI_APPLICATION))

	// Use default IDC_ARROW mouse cursor.
	wcex.Cursor = w32.LoadCursor(0, w32ext.MakeIntResource(w32.IDC_ARROW))

	// Assign HBRUSH to background using the standard window color
	wcex.Background = w32.COLOR_WINDOW + 11

	// Assign name of menu resource
	wcex.MenuName = nil

	wcex.ClassName = lpszClassName
	wcex.IconSm = w32.LoadIcon(hInstance, w32ext.MakeIntResource(w32.IDI_APPLICATION))

	return wcex
}