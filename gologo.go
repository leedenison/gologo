package main

//TODO: separate out the win32 stuff?
import "github.com/AllenDang/w32"

import (
	"C"
	"fmt"
	"github.com/leedenison/gologo/w32ext"
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

func UpdateStructures(wCtx *w32ext.WindowContext) {
	// Get the pane size
	winRect := w32ext.GetClientRect(wCtx)

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

func Tick(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
	wCtx.HDC = w32.GetDC(wCtx.Window)

	// TODO: Need to mutex this so we don't enter twice
	// Clear old ball
	w32ext.ClearRect(wCtx, ball.y, ball.x, 
		ball.y + ball.radius*2,
		ball.x + ball.radius*2)

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
			w32ext.KillTimer(wCtx, TIMER_ID)
			fmt.Printf("Hit wall %v\n", i)
		}
	}

	// Check if we've gone our right
	winRect := w32ext.GetClientRect(wCtx)
	if ball.y >= winRect.Right {
		w32ext.KillTimer(wCtx, TIMER_ID)
		if ball.x < 0 {
			fmt.Printf("Went over wall\n")
		} else {
			fmt.Printf("Win!\n")
		}
	}

	// Check if we've hit the left or bottom of the screen
	if ball.y+ball.radius*2 <= 0 || ball.x >= winRect.Bottom {
		w32ext.KillTimer(wCtx, TIMER_ID)
		fmt.Printf("Went out of play left or down\n")
	}

	PaintMovables(wCtx)
	w32ext.ReleaseDC(wCtx)
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

func PaintStructures(wCtx *w32ext.WindowContext) {
	for i := range walls {
		// Draw wall
		w32ext.DrawRectangle(
			wCtx,
			pens[PEN_WALL],
			walls[i].y,
			walls[i].x,
			walls[i].bottomy,
			walls[i].bottomx)
	}
}

func OnSize(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
	UpdateStructures(wCtx)
	OnPaint(wCtx, ev)
}

func OnPaint(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
	PaintStructures(wCtx)
	PaintMovables(wCtx)
}

func WinMain() int {
	aCtx := w32ext.GetAppContext()

	CreateWindowClass(&aCtx, "WNDclass")

	wCtx := CreateWindowInstance(&aCtx, "WNDclass", "Simple Go Window!")

	SetTimer(&wCtx, TIMER_ID, 100, Tick)
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
	EventHandlers[w32.WM_SIZE] = OnSize
	EventHandlers[w32.WM_PAINT] = OnPaint

	CreateObjects()
	WinMain()
	return
}
