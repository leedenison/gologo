package gologo

import (
    "C"
    "math"
    "github.com/AllenDang/w32"
)

const GRAVITY = 1
const RESISTANCE = 0
const MAX_SPEED = 10
const WALL_WIDTH = 10
const GAP_WIDTH = 80

const (
    OBJECT_SOLID = iota
    OBJECT_OUTLINE
    OBJECT_EMPTY
)

const (
    OVERLAP_FULL = iota
    OVERLAP_PARTIAL
)

const WINDOW_BORDER = 0

var next_object_id = 1
var objects = map[int]Object {}
var movables = map[int]Vector {}

type Vector struct {
    x, y   float64
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
    Radius float64
    Centre Vector
}

func (c *Circle) GetOrigin() Vector {
    return c.Centre
}

func (c *Circle) SetOrigin(v Vector) {
    c.Centre = v
}

//TODO set closest impact in return
func CheckCollisionPolyCirc(poly *Polygon, circ *Circle) *Collision {
    isRight := true

    for idx, v1 := range poly.Vertices {
        idx2 := (idx + 1) % len(poly.Vertices)
        v2 := poly.Vertices[idx2]

        // Check proximity to first vertex
        v1CircDist := Vector {x: v1.x - circ.Centre.x, y: v1.y - circ.Centre.y}
        objDistSq := math.Pow(v1CircDist.x, 2) + 
                 math.Pow(v1CircDist.y, 2)
        radiusSq := math.Pow(circ.Radius, 2)

        if objDistSq < radiusSq {
            // collision - the vertex is too close
            // (also catches v1 == v2)
            return &Collision{
                collisionType: OVERLAP_PARTIAL,
                closestImpact: Vector{ x: 0, y: 0 },
            }            
        }

        // Check proximity to edge
        v1v2Dist := Vector {x: v1.x - v2.x, y: v1.y - v2.y}
        edgeLengthSq := math.Pow(v1v2Dist.x, 2) + 
                     math.Pow(v1v2Dist.y, 2)
        dot := v1CircDist.x * v1v2Dist.x + v1CircDist.y * v1v2Dist.y
        normDist := Clamp(dot / edgeLengthSq, 0, 1) // Clamp as may be off edge end

        projectionX := v1.x + normDist * v1v2Dist.x
        projectionY := v1.y + normDist * v1v2Dist.y
        projDistSq := math.Pow(v1.x - projectionX, 2) +
                     math.Pow(v1.y - projectionY, 2)

        if projDistSq < radiusSq {
            return &Collision{
                collisionType: OVERLAP_PARTIAL,
                closestImpact: Vector{ x: 0, y: 0 },
            }                        
        }

        // Not had collision so record if it's on our right
        if (v1v2Dist.x * v1CircDist.y - v1v2Dist.y * v1CircDist.x) < 0 {
            isRight = false
        }
    }

    if isRight == true {
        return &Collision{
            collisionType: OVERLAP_FULL,
            closestImpact: Vector{ x: 0, y: 0 },
        }
    } else {
        return nil
    }
}

// TODO: Calulate if partial or full overlap and closest impact
func CheckCollisionPolyPoly(polyA *Polygon, polyB *Polygon) *Collision {
    for _, polygon := range []Polygon{ *polyA, *polyB } {
        for idx, v1 := range polygon.Vertices {
            idx2 := (idx + 1) % len(polygon.Vertices)
            v2 := polygon.Vertices[idx2]
            
            normal := Vector{v2.y - v1.y, v1.x - v2.x}

            var minA, maxA float64
            for _, p := range polyA.Vertices {
                projected := normal.x * p.x + normal.y * p.y;
                if minA == 0 || projected < minA { minA = projected }
                if maxA == 0 || projected > maxA { maxA = projected }
            }

            var minB, maxB float64
            for _, p := range polyB.Vertices {
                projected := normal.x * p.x + normal.y * p.y;
                if minB == 0 || projected < minB { minB = projected }
                if maxB == 0 || projected > maxB { maxB = projected }
            }

            if maxA < minB || maxB < minA { return nil }
        }
    }

    return &Collision {
        collisionType: OVERLAP_PARTIAL,
        closestImpact: Vector{ x: 0, y: 0 },
    }
}

func CheckCollisionCircCirc(circA *Circle, circB *Circle) *Collision {
        distX := circA.Centre.x - circB.Centre.x
        distY := circA.Centre.y - circB.Centre.y
        distSq := math.Pow(distX, 2) + math.Pow(distY, 2)
        radiusSq := math.Pow(circA.Radius + circB.Radius, 2)

        if distSq <= radiusSq {
            var closestImpact Vector
            
            if circA.Centre.x > circB.Centre.x {
                closestImpact.x = circB.Centre.x + distX / 2
            } else {
                closestImpact.x = circA.Centre.x + distX / 2
            }

            if circA.Centre.y > circB.Centre.y {
                closestImpact.y = circB.Centre.y + distY / 2
            } else {
                closestImpact.y = circA.Centre.y + distY / 2
            }

            halfRadiusASq := math.Pow(circA.Radius / 2, 2)
            halfRadiusBSq := math.Pow(circB.Radius / 2, 2)
            if distSq < halfRadiusASq ||
               distSq < halfRadiusBSq {
                return &Collision{
                    collisionType: OVERLAP_FULL,
                    closestImpact: closestImpact,
                }
            }

            return &Collision{
                collisionType: OVERLAP_PARTIAL,
                closestImpact: closestImpact,
            }

        }
        return nil
}

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
            Vector{ x: float64(clientRect.Right), y: 0 },
            Vector{ x: float64(clientRect.Right), y: float64(clientRect.Bottom) },
            Vector{ x: 0, y: float64(clientRect.Bottom) },
        },
        Renderer: renderers[RENDER_BG],
    }
}

func PhysicsTick(hwnd w32.HWND) {
    UpdateSpeedsAndPositions()
    // UpdateCollisions(wCtx)
}
