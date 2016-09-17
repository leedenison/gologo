package main

//TODO: separate out the win32 stuff?
import "github.com/AllenDang/w32"

import (
    "C"
    "fmt"
    "math"
    "reflect"
    "github.com/leedenison/gologo/w32ext"
)

const GOLOGO_MAIN_WIN = "GOLOGO_MAIN"

const SIXTY_HZ_IN_MILLIS = 16

const GRAVITY = 1
const RESISTANCE = 0
const MAX_SPEED = 10
const TIMER_ID = 1
const WALL_WIDTH = 10
const GAP_WIDTH = 80

const (
    OBJECT_SOLID = iota
    OBJECT_OUTLINE = iota
    OBJECT_EMPTY = iota
)

const (
    OVERLAP_FULL = iota
    OVERLAP_PARTIAL = iota
)

type Vector struct {
    x, y   int32
}

type Object interface {
    GetOrigin() Vector
    SetOrigin(Vector)
    GetResistance() int32
    RenderObjType() int
    CheckCollisionPolygon(*Polygon) *Collision
    CheckCollisionCircle(*Circle) *Collision
}

type Collision struct {
    collisionType int
    closestImpact Vector
}

// We currently assume all polygons are convex
type Polygon struct {
    vertices []Vector
    origin Vector
    resistance int32
    renderObjType int
    collisionObjType int
}

func (p *Polygon) GetOrigin() Vector {
    return p.origin
}

func (p *Polygon) SetOrigin(v Vector) {
    p.origin = v
}

func (p *Polygon) GetResistance() int32 {
    return p.resistance
}

func (p *Polygon) RenderObjType() int {
    return p.renderObjType
}

func (p *Polygon) CheckCollisionPolygon(target *Polygon) *Collision {
    return nil
}

func (p *Polygon) CheckCollisionCircle(target *Circle) *Collision {
    // TODO: If collisionObjType == OBJECT_SOLID check for fully inside
    var closestImpact Vector
    smallestDistance := math.Inf(1)

    for idx, v1 := range p.vertices {
    	idx2 := (idx + 1) % len(p.vertices)
        v2 := p.vertices[idx2]

        closestX := Clamp(target.center.x, v1.x, v2.x)
        closestY := Clamp(target.center.y, v1.y, v2.y)
        distanceX := target.center.x - closestX
        distanceY := target.center.y - closestY

        distanceSqrd := math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2)
        fmt.Printf("    Collision: Checking line (%v, %v)\n", v1, v2)
        fmt.Printf("        Collision: closestX = %v, closestY = %v\n", closestX, closestY)
        fmt.Printf("        Collision: distanceSqrd = %v, radiusSqrd = %v\n", distanceSqrd, math.Pow(float64(target.radius), 2))
        if distanceSqrd < math.Pow(float64(target.radius), 2) {
        	if math.IsInf(smallestDistance, 1) && distanceSqrd < smallestDistance {
                closestImpact = Vector{ x: closestX, y: closestY }
                fmt.Printf("        Collision: Recording closestX = %v, closestY = %v\n", closestX, closestY)
                smallestDistance = distanceSqrd        		
        	}
        }
    }

    if !math.IsInf(smallestDistance, 1) {
        return &Collision{ collisionType: OVERLAP_PARTIAL, closestImpact: closestImpact }        
    } else {
        return nil
    }
}

type Circle struct {
    radius, resistance int32
    center Vector
    renderObjType int
    collisionObjType int
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

func (c *Circle) RenderObjType() int {
    return c.renderObjType
}

func (c *Circle) CheckCollisionPolygon(target *Polygon) *Collision {
    return target.CheckCollisionCircle(c)
}

func (c *Circle) CheckCollisionCircle(target *Circle) *Collision {
    return nil
}

var objects = []Object {}
var speeds = map[int]Vector {}

const PEN_BALL = 0
const PEN_WALL = 1

var pens = map[int]*w32ext.Pen { 
    PEN_BALL: &w32ext.Pen{ Color: w32ext.RGB(0, 255, 0) },
    PEN_WALL: &w32ext.Pen{ Color: w32ext.RGB(0, 0, 0) },
}

func CreateObjects() {
    objects = append(objects, &Circle{
        center: Vector{ x: 50, y: 260 },
        radius: 20, 
        resistance: RESISTANCE,
        renderObjType: OBJECT_SOLID,
    })

    speeds[len(objects) - 1] = Vector{ x: 5, y: -5 }

    objects = append(objects, &Polygon{
        origin: Vector{ x: 0, y: 0 },
        vertices: []Vector{
            Vector{ x: 0, y: 0 }, 
            Vector{ x: 1024, y: 0 }, 
            Vector{ x: 1024, y: 768 }, 
            Vector{ x: 0, y: 768 }, 
        },
        resistance: RESISTANCE,
        renderObjType: OBJECT_EMPTY,
    })
}

func UpdateSpeedsAndPositions() {
    for objIdx, speed := range speeds {
        obj := objects[objIdx]
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

        speeds[objIdx] = speed

        origin := obj.GetOrigin()
        origin.x += speed.x
        origin.y += speed.y
        obj.SetOrigin(origin)
    }
}

func Tick(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
    // TODO: Need to mutex this so we don't enter twice
    w32ext.GetDC(wCtx)
    ClearMovables(wCtx)
    UpdateSpeedsAndPositions()
    UpdateCollisions(wCtx)
    PaintMovables(wCtx)
    w32ext.ReleaseDC(wCtx)
    w32ext.ReleaseContext(wCtx)
}

func UpdateCollisions(wCtx *w32ext.WindowContext) {
    for idx1, obj1 := range objects {
        for idx2 := idx1 + 1; idx2 < len(objects); idx2++ {
            fmt.Printf("Checking collisions - %v, %v\n", reflect.TypeOf(obj1), reflect.TypeOf(objects[idx2]))
            switch t2 := objects[idx2].(type) {
            case *Circle:
                collision := obj1.CheckCollisionCircle(t2)
                if collision != nil {
                    // hit a wall
                    w32ext.KillTimer(wCtx, TIMER_ID)
                    fmt.Printf("Hit wall %v\n", collision.closestImpact)
                }
            case *Polygon:
                collision := obj1.CheckCollisionPolygon(t2)
                if collision != nil {
                    // hit a wall
                    w32ext.KillTimer(wCtx, TIMER_ID)
                    fmt.Printf("Hit wall %v\n", collision.closestImpact)
                }
            }
        }
    }
}

func ClearMovables(wCtx *w32ext.WindowContext) {
    for _, obj := range objects {
        if obj.RenderObjType() != OBJECT_EMPTY {
            switch t := obj.(type) {
            case *Circle:
                w32ext.ClearRect(wCtx, 
                    t.center.x - t.radius,
                    t.center.y - t.radius, 
                    t.center.x + t.radius,
                    t.center.y + t.radius)
            default:
                _ = t
            }
        }
    }
}

func PaintMovables(wCtx *w32ext.WindowContext) {
    for _, obj := range objects {
        if obj.RenderObjType() != OBJECT_EMPTY {
            switch t := obj.(type) {
            case *Circle:
                w32ext.DrawEllipse(
                    wCtx,
                    pens[PEN_BALL], 
                    t.center.x - t.radius, 
                    t.center.y - t.radius, 
                    t.center.x + t.radius,
                    t.center.y + t.radius)
            default:
                _ = t
            }            
        }
    }
}

func OnSize(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
    OnPaint(wCtx, ev)
}

func OnPaint(wCtx *w32ext.WindowContext, ev *w32ext.Event) {
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
