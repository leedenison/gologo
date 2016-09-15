package main

//TODO: separate out the win32 stuff?
import "github.com/AllenDang/w32"

import (
	"C"
	"fmt"
	"github.com/leedenison/gologo/gologowin"
	"github.com/leedenison/gologo/w32ext"
	"syscall"
)

const GRAVITY_VALUE = 5
const CIRCLE_RESISTANCE = 3
const MAX_SPEED = 40
const TIMER_ID = 1
const WALL_WIDTH = 10
const GAP_WIDTH = 80

type shape struct {
	x, y   int32
}

type structure struct {
	bottomx, bottomy int32
	shape
}

type movableshape struct {
	vx, vy int32
	shape
}

// Feels like the Resistable interface and the movableshape
// struct are actually the same thing but you can't combine them
// Also this interface isn't actually used yet but theoretically
// it could be if we had more than one movable object
type Resistable interface {
	ApplyRst()
}

type circle struct {
	radius, resistance int32
	movableshape
}

func (s *movableshape) Move() {
	s.x += s.vx
	s.y += s.vy
}

func (s *movableshape) ApplyGrav() {
	s.vx += GRAVITY_VALUE
	if s.vx > MAX_SPEED {
		s.vx = MAX_SPEED
	}
}

func (c *circle) ApplyRst() {
	if c.vy > 0 {
		c.vy -= c.resistance
		if c.vy < 0 {
			c.vy = 0
		}
	} else if c.vy < 0 {
		c.vy += c.resistance
		if c.vy > 0 {
			c.vy = 0
		}
	}
}

var ball = circle{}

// TODO Don't know go conventions
// - structures for the var name seems more sensible but then
// seems to be trying hard to be confusing
var walls []structure

const PEN_BALL = 0
const PEN_WALL = 1

var pens = map[int]*w32ext.Pen { 
	PEN_BALL: &w32ext.Pen{ Color: w32ext.RGB(0, 255, 0) },
	PEN_WALL: &w32ext.Pen{ Color: w32ext.RGB(0, 0, 1) },
}

func Clamp(value, min, max int32) int32 {
	switch {
	case value < min:
		return min
	case value > max:
		return max
	}
	return value
}

func CreateObjects() {
	ball = circle{
		movableshape: movableshape{
			shape: shape{ x: 240, y: 30 },
			vx: -35, 
			vy: 50,
		},
		radius: 20, 
		resistance: CIRCLE_RESISTANCE,
	}

	for i := 0; i < 4; i++ {
		wall := structure{}
		walls = append(walls, wall)
	}
}

func UpdateStructures(hwnd w32.HWND) {
	// Get the pane size
	winRect := w32.GetClientRect(hwnd)

	// Do the left wall
	walls[0].bottomx = winRect.Bottom
	walls[0].bottomy = WALL_WIDTH
	// Floor
	walls[1].x = winRect.Bottom - WALL_WIDTH
	walls[1].bottomx = winRect.Bottom
	walls[1].bottomy = winRect.Right
	// Top right wall
	walls[2].y = winRect.Right - WALL_WIDTH
	walls[2].bottomx = winRect.Bottom/2 - GAP_WIDTH/2
	walls[2].bottomy = winRect.Right
	// Bottom right wall
	walls[3].x = winRect.Bottom/2 + GAP_WIDTH/2
	walls[3].y = winRect.Right - WALL_WIDTH
	walls[3].bottomx = winRect.Bottom
	walls[3].bottomy = winRect.Right
}

func Tick(hwnd w32.HWND, uMsg uint32, idEvent uintptr,
	dwTime uint16) uintptr {
	// TODO: Need to mutex this so we don't enter twice
	// Get ball rect
	ballRect := w32.RECT{Left: ball.y, Top: ball.x,
		Bottom: ball.x + ball.radius*2,
		Right:  ball.y + ball.radius*2}
	// Clear ball
	w32.InvalidateRect(hwnd, &ballRect, true)

	// Get balls new position
	ball.ApplyGrav()
	ball.ApplyRst()
	ball.Move()

	// Check for collisions with walls
	for i := range walls {
		closestX := Clamp(ball.x+ball.radius, walls[i].x, walls[i].bottomx)
		closestY := Clamp(ball.y+ball.radius, walls[i].y, walls[i].bottomy)
		distanceX := ball.x + ball.radius - closestX
		distanceY := ball.y + ball.radius - closestY

		distanceSqrd := distanceX*distanceX + distanceY*distanceY
		if distanceSqrd < ball.radius*ball.radius {
			// hit a wall
			w32ext.KillTimer(hwnd, uintptr(TIMER_ID))
			fmt.Printf("Hit wall %v\n", i)
		}
	}

	// Check if we've gone our right
	winRect := w32.GetClientRect(hwnd)
	if ball.y >= winRect.Right {
		w32ext.KillTimer(hwnd, uintptr(TIMER_ID))
		if ball.x < 0 {
			fmt.Printf("Went over wall\n")
		} else {
			fmt.Printf("Win!\n")
		}
	}

	// Check if we've hit the left or bottom of the screen
	if ball.y+ball.radius*2 <= 0 ||
		ball.x >= winRect.Bottom {
		w32ext.KillTimer(hwnd, uintptr(TIMER_ID))
		fmt.Printf("Went out of play left or down\n")
	}

	hdc := w32.GetDC(hwnd)
	wCtx := &w32ext.WindowContext { Window: hwnd, HDC: hdc }
	PaintMovables(wCtx)
	w32.ReleaseDC(hwnd, hdc)

	return 0
}

func PaintMovables(wCtx *w32ext.WindowContext) {
	// Draw ball
	w32ext.DrawEllipse(
		wCtx,
		pens[PEN_BALL], 
		ball.y, 
		ball.x, 
		ball.y+ball.radius*2,
		ball.x+ball.radius*2)
}

func PaintStructures(hdc w32.HDC) {

	if len(walls) == 0 {
		return
	}
	// Create previous brush and pen holder
	var previousPen, previousBrush w32.HGDIOBJ

	for i := range walls {
		// Create brush and pen for this wall
		lBrush := w32.LOGBRUSH{LbStyle: w32.BS_SOLID, LbColor: w32.COLORREF(pens[PEN_WALL].Color)}
		hBrush := w32.CreateBrushIndirect(&lBrush)
		hPen := w32.ExtCreatePen(w32.PS_COSMETIC|w32.PS_SOLID, 1, &lBrush, 0, nil)

		// Select pen
		oldPen := w32.SelectObject(hdc, w32.HGDIOBJ(hPen))
		oldBrush := w32.SelectObject(hdc, w32.HGDIOBJ(hBrush))

		// TODO Seems hideous to do this but don't know what's better
		if i == 0 {
			previousPen = oldPen
			previousBrush = oldBrush
		}

		// Draw wall
		w32.Rectangle(hdc, int(walls[i].y), int(walls[i].x),
			int(walls[i].bottomy), int(walls[i].bottomx))

		// Delete objects
		w32.DeleteObject(w32.HGDIOBJ(hPen))
		w32.DeleteObject(w32.HGDIOBJ(hBrush))
	}

	// Reselect previous pen
	w32.SelectObject(hdc, previousPen)
	w32.SelectObject(hdc, previousBrush)
}

func WndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case w32.WM_DESTROY:
		// 0 = WM_QUIT
		w32.PostQuitMessage(0)
	case w32.WM_SIZE:
		UpdateStructures(hwnd)
		fallthrough
	case w32.WM_PAINT:
		// On initial paint
		var ps w32.PAINTSTRUCT

		hdc := w32.BeginPaint(hwnd, &ps)
		wCtx := &w32ext.WindowContext { Window: hwnd, HDC: hdc }
		PaintStructures(hdc)
		PaintMovables(wCtx)
		w32.EndPaint(hwnd, &ps)
	default:
		return w32.DefWindowProc(hwnd, msg, wParam, lParam)
	}
	return 0
}

func WinMain() int {

	// Handle to application instance.
	hInstance := w32.GetModuleHandle("")

	// Registered class name of the window.
	lpszClassName := syscall.StringToUTF16Ptr("WNDclass")

	wcex := gologowin.CreateWindowClass(hInstance, lpszClassName, syscall.NewCallback(WndProc))

	// Make this window available to other controls
	w32.RegisterClassEx(&wcex)

	// Create an instance of this window class
	hwnd := w32.CreateWindowEx(0, lpszClassName,
		syscall.StringToUTF16Ptr("Simple Go Window!"),
		w32.WS_OVERLAPPEDWINDOW|w32.WS_VISIBLE,
		w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, 400, 400, 0, 0,
		hInstance, nil)

	w32.ShowWindow(hwnd, w32.SW_SHOWDEFAULT)
	w32.UpdateWindow(hwnd)

	w32ext.SetTimer(hwnd, uintptr(TIMER_ID), 100, syscall.NewCallback(Tick))
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
	CreateObjects()
	WinMain()
	return
}
