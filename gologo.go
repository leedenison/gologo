package main

import "github.com/AllenDang/w32"

import (
	"C"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")

    procSetTimer = moduser32.NewProc("SetTimer")
    procKillTimer = moduser32.NewProc("KillTimer")
)

func SetTimer(hwnd w32.HWND, nIDEvent uintptr, uElapse uint32, lpTimerProc uintptr) uintptr {
	ret, _, _ := procSetTimer.Call(
		uintptr(hwnd),
		nIDEvent,
		uintptr(uElapse),
		lpTimerProc)

	return ret
}

func KillTimer(hwnd w32.HWND, nIDEvent uintptr) bool {
	ret, _, _ := procKillTimer.Call(
		uintptr(hwnd),
		nIDEvent)

	return ret != 0
}

func MakeIntResource(id uint16) (*uint16) {
    return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func RGB(r, g, b byte) (uint32) {
	return (uint32(r) | uint32(g) << 8 | uint32(b) << 16)
}

func Tick(hwnd w32.HWND, uMsg uint32, idEvent uintptr, dwTime uint16) (uintptr) {
	fmt.Printf("Tick!\n")

	return 0
}

func OnPaint(hdc w32.HDC) {
	green := RGB(0, 255, 0)

	// Create brush
    lBrush := w32.LOGBRUSH{LbStyle: w32.BS_SOLID, LbColor: w32.COLORREF(green)}
    hBrush := w32.CreateBrushIndirect(&lBrush)

    // Create pen
    hPen := w32.ExtCreatePen(w32.PS_COSMETIC|w32.PS_SOLID, 1, &lBrush, 0, nil)

    // Select pen
    previousPen := w32.SelectObject(hdc, w32.HGDIOBJ(hPen))

    // Draw line
    w32.MoveToEx(hdc, 0, 0, nil)
    w32.LineTo(hdc, 200, 200)

    // Reselect previous pen
    w32.SelectObject(hdc, previousPen)

    // Delete objects
    w32.DeleteObject(w32.HGDIOBJ(hPen))
    w32.DeleteObject(w32.HGDIOBJ(hBrush))
}

func WndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) (uintptr) {
	switch msg {
	case w32.WM_DESTROY:
		// 0 = WM_QUIT
        w32.PostQuitMessage(0)
    case w32.WM_PAINT:
    	var ps w32.PAINTSTRUCT

    	hdc := w32.BeginPaint(hwnd, &ps)
        OnPaint(hdc)
        w32.EndPaint(hwnd, &ps)
    default:
        return w32.DefWindowProc(hwnd, msg, wParam, lParam)
    }
    return 0
}

func CreateWindowClass(hInstance w32.HINSTANCE, lpszClassName *uint16) (w32.WNDCLASSEX) {
	var wcex w32.WNDCLASSEX

	// Size of the window object.
	wcex.Size = uint32(unsafe.Sizeof(wcex))

	wcex.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
	// Application loop handler procedure.
	wcex.WndProc = syscall.NewCallback(WndProc)
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
	wcex.Icon = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))

	// Use default IDC_ARROW mouse cursor.
	wcex.Cursor = w32.LoadCursor(0, MakeIntResource(w32.IDC_ARROW))

	// Assign HBRUSH to background using the standard window color
	wcex.Background = w32.COLOR_WINDOW + 11

	// Assign name of menu resource
	wcex.MenuName = nil

	wcex.ClassName = lpszClassName
	wcex.IconSm = w32.LoadIcon(hInstance, MakeIntResource(w32.IDI_APPLICATION))

	return wcex
}

const TIMER_ID = 1

func WinMain() int {

	// Handle to application instance.
	hInstance := w32.GetModuleHandle("")

	// Registered class name of the window.
	lpszClassName := syscall.StringToUTF16Ptr("WNDclass")

	wcex := CreateWindowClass(hInstance, lpszClassName)

	// Make this window available to other controls
	w32.RegisterClassEx(&wcex)

	// Create an instance of this window class
	hwnd := w32.CreateWindowEx(0, lpszClassName, syscall.StringToUTF16Ptr("Simple Go Window!"), w32.WS_OVERLAPPEDWINDOW | w32.WS_VISIBLE, w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 400, 0, 0, hInstance, nil)

	w32.ShowWindow(hwnd, w32.SW_SHOWDEFAULT)
	w32.UpdateWindow(hwnd)

	SetTimer(hwnd, uintptr(TIMER_ID), 1000, syscall.NewCallback(Tick))
   	var msg w32.MSG
   	for {
   		// 0, 0, 0 = retrive all messages from all sources
        if w32.GetMessage(&msg, 0, 0, 0) == 0 {
            break
        }
        w32.TranslateMessage(&msg)
        w32.DispatchMessage(&msg)
   	}
   	return int(msg.WParam)
}

func main() {
    WinMain()
    return
}
