package main

//TODO: separate out the win32 stuff?
import "github.com/AllenDang/w32"

import (
	"C"
	"fmt"
	"math"
	"github.com/leedenison/gologo/w32ext"
)

const GOLOGO_MAIN_WIN = "GOLOGO_MAIN"

const SIXTY_HZ_IN_MILLIS = 16

const GRAVITY = 5
const RESISTANCE = 3
const MAX_SPEED = 40
const TIMER_ID = 1
const WALL_WIDTH = 10
const GAP_WIDTH = 80

type Vector struct {
	x, y   int32
}

type Moveable interface {
	GetOrigin() Vector
	SetOrigin(Vector)
	GetResistance() int32
}

type Rectangle struct {
	topleft, bottomright Vector
}

type Circle struct {
	radius, resistance int32
	center Vector
}

func (c *Circle) GetOrigin() Vector {
	return c.center
}

func (c *Circle) SetOrigin(v Vector) {
	c.center = v
}

func (c *Circle) GetResistance() int32 {
	return c.resistance
}

var ball Circle

// TODO Don't know go conventions
// - structures for the var name seems more sensible but then
// seems to be trying hard to be confusing
var walls []Rectangle

var speeds = map[Moveable]Vector {}

const PEN_BALL = 0
const PEN_WALL = 1

var pens = map[int]*w32ext.Pen { 
	PEN_BALL: &w32ext.Pen{ Color: w32ext.RGB(0, 255, 0) },
	PEN_WALL: &w32ext.Pen{ Color: w32ext.RGB(0, 0, 1) },
}

func CreateObjects() {
	ball = Circle{
		center: Vector{ x: 50, y: 260 },
		radius: 20, 
		resistance: RESISTANCE,
	}

	speeds[&ball] = Vector{ x: 50, y: -35 }

	for i := 0; i < 4; i++ {
		walls = append(walls, Rectangle{})
	}
}

func UpdateStructures(wCtx *w32ext.WindowContext) {
	// Get the pane size
	winRect := w32ext.GetClientRect(wCtx)

	// Do the left wall
	walls[0].bottomright.x = WALL_WIDTH
	walls[0].bottomright.y = winRect.Bottom
	// Floor
	walls[1].topleft.y = winRect.Bottom - WALL_WIDTH
	walls[1].bottomright.x = winRect.Right
	walls[1].bottomright.y = winRect.Bottom
	// Top right wall
	walls[2].topleft.x = winRect.Right - WALL_WIDTH
	walls[2].bottomright.x = winRect.Right
	walls[2].bottomright.y = winRect.Bottom/2 - GAP_WIDTH/2
	// Bottom right wall
	walls[3].topleft.x = winRect.Right - WALL_WIDTH
	walls[3].topleft.y = winRect.Bottom/2 + GAP_WIDTH/2
	walls[3].bottomright.x = winRect.Right
	walls[3].bottomright.y = winRect.Bottom
}

func UpdateSpeedsAndPositions() {
    for obj, speed := range speeds {
       	speed.y += GRAVITY
	 	if speed.y > MAX_SPEED {
			speed.y = MAX_SPEED
		}

		if speed.x > 0 {
			speed.x -= obj.GetResistance()
			if speed.x < 0 {
				speed.x = 0
			}
		} else if speed.x < 0 {
			speed.x += obj.GetResistance()
			if speed.x > 0 {
				speed.x = 0
			}
		}

		speeds[obj] = speed

		origin := obj.GetOrigin()
		origin.x += speed.x
		origin.y += speed.y
		obj.SetOrigin(origin)
	}
}

func Tick(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
	w32ext.GetDC(wCtx)

	// TODO: Need to mutex this so we don't enter twice
	// Clear old ball
	w32ext.ClearRect(wCtx, 
		ball.center.x - ball.radius,
		ball.center.y - ball.radius, 
		ball.center.x + ball.radius,
		ball.center.y + ball.radius)

	// Get balls new position
	UpdateSpeedsAndPositions()

	// Check for collisions with walls
	for i := range walls {
		closestX := Clamp(ball.center.x, walls[i].topleft.x, walls[i].bottomright.x)
		closestY := Clamp(ball.center.y, walls[i].topleft.y, walls[i].bottomright.y)
		distanceX := ball.center.x - closestX
		distanceY := ball.center.y - closestY

		distanceSqrd := math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2)
		if distanceSqrd < math.Pow(float64(ball.radius), 2) {
			// hit a wall
			w32ext.KillTimer(wCtx, TIMER_ID)
			fmt.Printf("Hit wall %v\n", i)
		}
	}

	// Check if we've gone out right
	winRect := w32ext.GetClientRect(wCtx)
	if ball.center.x - ball.radius >= winRect.Right {
		w32ext.KillTimer(wCtx, TIMER_ID)
		if ball.center.y < 0 {
			fmt.Printf("Went over wall\n")
		} else {
			fmt.Printf("Win!\n")
		}
	}

	// Check if we've hit the left or bottom of the screen
	if ball.center.x+ball.radius <= 0 || ball.center.y >= winRect.Bottom {
		w32ext.KillTimer(wCtx, TIMER_ID)
		fmt.Printf("Went out of play left or down\n")
	}

	PaintMovables(wCtx)
	w32ext.ReleaseDC(wCtx)
	w32ext.ReleaseContext(wCtx)
}

func PaintMovables(wCtx *w32ext.WindowContext) {
	// Draw ball
	w32ext.DrawEllipse(
		wCtx,
		pens[PEN_BALL], 
		ball.center.x - ball.radius, 
		ball.center.y - ball.radius, 
		ball.center.x + ball.radius,
		ball.center.y + ball.radius)
}

func PaintStructures(wCtx *w32ext.WindowContext) {
	for i := range walls {
		// Draw wall
		w32ext.DrawRectangle(
			wCtx,
			pens[PEN_WALL],
			walls[i].topleft.x,
			walls[i].topleft.y,
			walls[i].bottomright.x,
			walls[i].bottomright.y)
	}
}

func OnSize(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
	UpdateStructures(wCtx)
	OnPaint(wCtx, ev)
}

func OnPaint(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
	PaintStructures(wCtx)
	PaintMovables(wCtx)
	w32ext.ReleaseContext(wCtx)
}

func main() {
	aCtx := w32ext.GetAppContext()

	EventHandlers[w32.WM_SIZE] = OnSize
	EventHandlers[w32.WM_PAINT] = OnPaint

	CreateObjects()

	CreateWindowClass(&aCtx, GOLOGO_MAIN_WIN)
	wCtx := CreateWindowInstance(&aCtx, GOLOGO_MAIN_WIN, "Simple Go Window!")

	SetTimer(&wCtx, TIMER_ID, SIXTY_HZ_IN_MILLIS, Tick)

	var msg w32.MSG
	for {
		// 0, 0, 0 = retrive all messages from all sources
		if w32.GetMessage(&msg, 0, 0, 0) == 0 {
			break
		}
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}

	return
}
