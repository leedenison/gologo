package gologo

import (
    "C"
    "github.com/AllenDang/w32"
)

const GRAVITY = 1
const RESISTANCE = 0
const MAX_SPEED = 10
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

const WINDOW_BORDER = 0

var next_object_id = 1
var objects = map[int]Object {}
var movables = map[int]Vector {}

type Vector struct {
    x, y   int32
}

type Object interface {
    GetOrigin() Vector
    SetOrigin(Vector)
}

type Collision struct {
    collisionType int
    closestImpact Vector
}

// We currently assume all polygons are convex
type Polygon struct {
    Renderer *Renderer
    Vertices []Vector
    Origin Vector
}

func (p *Polygon) GetOrigin() Vector {
    return p.Origin
}

func (p *Polygon) SetOrigin(v Vector) {
    p.Origin = v
}

type Circle struct {
    Renderer *Renderer
    Radius int32
    Center Vector
}

func (c *Circle) GetOrigin() Vector {
    return c.Center
}

func (c *Circle) SetOrigin(v Vector) {
    c.Center = v
}

/*
func CheckCollisionCircle(p *Polygon, target *Circle) *Collision {
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

        distanceSqrd := math.Pow(float64(distanceX), 2) +
                math.Pow(float64(distanceY), 2)
        fmt.Printf("    Collision: Checking line (%v, %v)\n", v1, v2)
        fmt.Printf("        Collision: closestX = %v, closestY = %v\n",
                closestX,
                closestY)
        fmt.Printf("        Collision: distanceSqrd = %v, radiusSqrd = %v\n",
                distanceSqrd,
                math.Pow(float64(target.radius), 2))
        if distanceSqrd < math.Pow(float64(target.radius), 2) {
            if math.IsInf(smallestDistance, 1) &&
                    distanceSqrd < smallestDistance {
                closestImpact = Vector{ x: closestX, y: closestY }
                fmt.Printf(
                    "        Collision: Recording closestX = %v," +
                    " closestY = %v\n",
                    closestX,
                    closestY)
                smallestDistance = distanceSqrd
            }
        }
    }

    if !math.IsInf(smallestDistance, 1) {
        return &Collision{
            collisionType: OVERLAP_PARTIAL,
            closestImpact: closestImpact,
        }
    } else {
        return nil
    }
}
*/

func UpdateSpeedsAndPositions() {
    for objIdx, speed := range movables {
        obj := objects[objIdx]
        speed.y += GRAVITY

        movables[objIdx] = speed

        origin := obj.GetOrigin()
        origin.x += speed.x
        origin.y += speed.y
        obj.SetOrigin(origin)
    }
}

/*
func UpdateCollisions(wCtx *w32ext.WindowContext) {
    for idx1, obj1 := range objects {
        for idx2 := idx1 + 1; idx2 < len(objects); idx2++ {
            fmt.Printf(
                "Checking collisions - %v, %v\n",
                reflect.TypeOf(obj1),
                reflect.TypeOf(objects[idx2]))
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
*/

func UpdateWindowEdge(hwnd w32.HWND) {
    clientRect := w32.GetClientRect(hwnd)

    objects[WINDOW_BORDER] = &Polygon{
        Origin: Vector{ x: 0, y: 0 },
        Vertices: []Vector{
            Vector{ x: 0, y: 0 },
            Vector{ x: clientRect.Right, y: 0 },
            Vector{ x: clientRect.Right, y: clientRect.Bottom },
            Vector{ x: 0, y: clientRect.Bottom },
        },
        Renderer: renderers[RENDER_BG],
    }
}

func PhysicsTick(hwnd w32.HWND) {
    UpdateSpeedsAndPositions()
    // UpdateCollisions(wCtx)
}
