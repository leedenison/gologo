package main

import "github.com/AllenDang/w32"

import (
	"C"
	"syscall"
	"unsafe"
	//"fmt"
	"github.com/leedenison/gologo/timer"
)

const GRAVITY_VALUE = 5
const CIRCLE_RESISTANCE = 3
const MAX_SPEED = 40
const TIMER_ID = 1

type shape struct {
	x, y int
	colour uint32
}

type movableshape struct {
	vx, vy int
	shape
}

type Resistable interface {
	ApplyRst()
}

type circle struct {
	radius, resistance int
	movableshape
}

func (s *movableshape) Move() {
	s.x += s.vx
	s.y += s.vy
}

func (s *movableshape) ApplyGrav() {
	s.vy += GRAVITY_VALUE
	if (s.vy > MAX_SPEED) {s.vy = MAX_SPEED}
}

func (c *circle) ApplyRst() {
	c.vx -= c.resistance
	if (c.vx < 0) {c.vx = 0}
}

var ball = circle{movableshape: movableshape{shape: shape{x: 30, y: 240, colour: RGB(0,255,0)}, vx: 40, vy: -50}, radius: 20, resistance: CIRCLE_RESISTANCE}

func MakeIntResource(id uint16) (*uint16) {
    return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func RGB(r, g, b byte) (uint32) {
	return (uint32(r) | uint32(g) << 8 | uint32(b) << 16)
}

func Tick(hwnd w32.HWND, uMsg uint32, idEvent uintptr, dwTime uint16) (uintptr) {
	// TODO: Need to mutex this so we don't enter twice
	// Get ball rect
	ballRect := w32.RECT{Left: int32(ball.x), Top: int32(ball.y),
						Bottom: int32(ball.y + ball.radius * 2), 
						Right: int32(ball.x + ball.radius * 2)}
	// Clear ball
	w32.InvalidateRect(hwnd, &ballRect, true)

	// Get balls new position
	ball.ApplyGrav()
	ball.ApplyRst()
	ball.Move()

	// Check if we've hit the edge of the screen
	// Don't worry about top - we could come back on
	winRect := w32.GetClientRect(hwnd)
	if (ball.x <= 0 ||
		int32(ball.x + ball.radius * 2) >= winRect.Right ||
		int32(ball.y + ball.radius * 2) >= winRect.Bottom) {
		timer.KillTimer(hwnd, uintptr(TIMER_ID))
	}

	hdc := w32.GetDC(hwnd)
	OnPaint(hdc)
	w32.ReleaseDC(hwnd, hdc)
	
	return 0
}

func OnPaint(hdc w32.HDC) {
	// Create brush
    lBrush := w32.LOGBRUSH{LbStyle: w32.BS_SOLID, LbColor: w32.COLORREF(ball.colour)}
    hBrush := w32.CreateBrushIndirect(&lBrush)

    // Create pen
    hPen := w32.ExtCreatePen(w32.PS_COSMETIC|w32.PS_SOLID, 1, &lBrush, 0, nil)

    // Select pen
    previousPen := w32.SelectObject(hdc, w32.HGDIOBJ(hPen))

    // Draw ball
    w32.Ellipse(hdc, ball.x, ball.y, ball.x + ball.radius * 2, ball.y + ball.radius * 2)

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

	timer.SetTimer(hwnd, uintptr(TIMER_ID), 100, syscall.NewCallback(Tick))
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
