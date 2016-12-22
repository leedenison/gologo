package gologo

import (
    "C"
    "math"
    "github.com/AllenDang/w32"
)

const GRAVITY = 0.5
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

const WINDOW_BORDER_TOP = 0
const WINDOW_BORDER_LEFT = 1
const WINDOW_BORDER_BOTTOM = 2
const WINDOW_BORDER_RIGHT = 3

var next_object_id = 4
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
    Id int
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

func (p *Polygon) GetWorldVertices() []Vector {
    result := []Vector {}

    for _, v := range p.Vertices {
        result = append(result, Vector{
            x: v.x + p.Origin.x,
            y: v.y + p.Origin.y,
        })
    }

    return result
}

type Circle struct {
    Id int
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
    //fmt.Printf("Checking collision between:\n\tPolygon: %v\n\tCircle: %v\n", *poly, *circ)
    isRight := true
    polyVertices := poly.GetWorldVertices()

    //fmt.Printf("Looping over edges:\n")
    for idx, v1 := range polyVertices {
        idx2 := (idx + 1) % len(polyVertices)
        v2 := polyVertices[idx2]
        //fmt.Printf("\tTesting edge v1: %v, v2: %v\n", v1, v2)

        // leedenison: This looks like it is the vector from the circle to the
        // start point but I think it should be the vector from the start point
        // to the circle
        // Check proximity to first vertex
        //v1CircDist := Vector {x: v1.x - circ.Centre.x, y: v1.y - circ.Centre.y}
        v1CircDist := Vector {x: circ.Centre.x - v1.x, y: circ.Centre.y - v1.y}        
        //fmt.Printf("\t\tv1CircDirVec: %v\n", v1CircDist)
        radiusSq := math.Pow(circ.Radius, 2)
        //fmt.Printf("\t\tcirc.Radius ^ 2: %v\n", radiusSq)

        // Not sure why we need this independent check since the check below
        // is clamped to the length of the line
        /*
        objDistSq := math.Pow(v1CircDist.x, 2) + math.Pow(v1CircDist.y, 2)
        if objDistSq < radiusSq {
            // collision - the vertex is too close
            // (also catches v1 == v2)
            return &Collision{
                collisionType: OVERLAP_PARTIAL,
                closestImpact: Vector{ x: 0, y: 0 },
            }            
        }
        */

        // leedenison: This also looks reversed to me.
        // Check proximity to edge
        // v1v2Dist := Vector {x: v1.x - v2.x, y: v1.y - v2.y}
        v1v2Dist := Vector {x: v2.x - v1.x, y: v2.y - v1.y}

        //fmt.Printf("\t\tv1v2Vec: %v\n", v1v2Dist)
        edgeLengthSq := math.Pow(v1v2Dist.x, 2) + math.Pow(v1v2Dist.y, 2)
        //fmt.Printf("\t\t|v1v2Vec| ^ 2: %v\n", edgeLengthSq)

        dot := v1CircDist.x * v1v2Dist.x + v1CircDist.y * v1v2Dist.y
        //fmt.Printf("\t\tv1CircVec.v1v2Vec: %v\n", dot)
        normDist := Clamp(dot / edgeLengthSq, 0, 1) // Clamp as may be off edge end
        //fmt.Printf("\t\tDistance to nearest approach along edge (from v1): %v\n", normDist)        

        projectionX := v1.x + normDist * v1v2Dist.x
        projectionY := v1.y + normDist * v1v2Dist.y
        //fmt.Printf("\t\tnearApproachVec: %v, %v\n", projectionX, projectionY)        

        // leedenison: This seems to be the distance to the first vertex of the
        // polygon which seems wrong, I think it should be to the center of the circle
        //projDistSq := math.Pow(v1.x - projectionX, 2) +
        //             math.Pow(v1.y - projectionY, 2)
        projDistSq := math.Pow(circ.Centre.x - projectionX, 2) +
                     math.Pow(circ.Centre.y - projectionY, 2)
        //fmt.Printf("\t\t|nearApproachVec| ^ 2: %v\n", projDistSq)

        // We return the first collision we find, we may want to find the closest
        // collision in future
        //fmt.Printf("\t\tTest if |nearApproachVec| ^ 2 < circ.Radius ^ 2: %v\n", projDistSq < radiusSq)        
        if projDistSq < radiusSq {
            // fmt.Printf("Collision found returning: %v\n", Collision{
            //     collisionType: OVERLAP_PARTIAL,
            //     closestImpact: Vector{ x: projectionX, y: projectionY },
            // })
            return &Collision{
                collisionType: OVERLAP_PARTIAL,
                closestImpact: Vector{ x: projectionX, y: projectionY },
            }                        
        }

        //fmt.Printf("\t\tIs circle on the right of the edge: %v\n", (v1v2Dist.x * v1CircDist.y - v1v2Dist.y * v1CircDist.x) < 0)
        // Not had collision so record if it's on our right
        if (v1v2Dist.x * v1CircDist.y - v1v2Dist.y * v1CircDist.x) < 0 {
            isRight = false
        }
        //fmt.Printf("\t\tNo collision found, next edge.\n")
    }

    if isRight == true {
        // fmt.Printf("No more edges, but circle is on the right of all edges, so return: %v\n", Collision{
        //     collisionType: OVERLAP_FULL,
        //     closestImpact: Vector{ x: 0, y: 0 },
        // })
        return &Collision{
            collisionType: OVERLAP_FULL,
            closestImpact: Vector{ x: 0, y: 0 },
        }
    } else {
        //fmt.Printf("No collision found, returning <nil>\n")
        return nil
    }
}

func HandleCollisionPolyCirc(poly *Polygon, circ *Circle, collision *Collision) {
    // Current approach assumes the polygon has infinite mass (ie. cannot be moved)
    // and is stationary (ie. does not add energy to the circle).

    if velocity, ok := movables[circ.Id]; ok {
        // Calculate unit direction vector closestImpact -> centre
        dirVector := UnitDirectionVector(circ.Centre, collision.closestImpact)
        dot := DotProduct(velocity, dirVector)

        // Reflect the circle velocity around the direction vector
        velocity.x = velocity.x - 2 * dot * dirVector.x
        velocity.y = velocity.y - 2 * dot * dirVector.y

        movables[circ.Id] = velocity
    }
}

// TODO: Calulate if partial or full overlap and closest impact
func CheckCollisionPolyPoly(polyA *Polygon, polyB *Polygon) *Collision {
    polyAVertices := polyA.GetWorldVertices()
    polyBVertices := polyB.GetWorldVertices()

    for _, polyVertices := range [][]Vector{ polyAVertices, polyBVertices } {
        for idx, v1 := range polyVertices {
            idx2 := (idx + 1) % len(polyVertices)
            v2 := polyVertices[idx2]
            
            normal := Vector{v2.y - v1.y, v1.x - v2.x}

            var minA, maxA float64
            for _, p := range polyAVertices {
                projected := normal.x * p.x + normal.y * p.y;
                if minA == 0 || projected < minA { minA = projected }
                if maxA == 0 || projected > maxA { maxA = projected }
            }

            var minB, maxB float64
            for _, p := range polyBVertices {
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

func UpdateCollisions() {
    for idx1, obj1 := range objects {
        for idx2 := idx1 + 1; idx2 < len(objects); idx2++ {
            switch t1 := obj1.(type) {
            case *Circle:
                switch t2 := objects[idx2].(type) {
                case *Polygon:
                    collision := CheckCollisionPolyCirc(t2, t1)
                    if collision != nil {
                        HandleCollisionPolyCirc(t2, t1, collision)
                    }
                }
            case *Polygon:
                switch t2 := objects[idx2].(type) {
                case *Circle:
                    collision := CheckCollisionPolyCirc(t1, t2)
                    if collision != nil {
                        HandleCollisionPolyCirc(t1, t2, collision)
                    }
                }
            }
        }
    }
}

func UpdateWindowEdge(hwnd w32.HWND) {
    clientRect := w32.GetClientRect(hwnd)

    objects[WINDOW_BORDER_TOP] = &Polygon{
        Id: WINDOW_BORDER_TOP,
        Origin: Vector{ x: 0, y: 0 },
        Vertices: []Vector{
            Vector{ x: 0, y: 0 },
            Vector{ x: 0, y: -1 },
            Vector{ x: float64(clientRect.Right), y: -1 },
            Vector{ x: float64(clientRect.Right), y: 0 },
        },
        Renderer: renderers[RENDER_BG],
    }
    objects[WINDOW_BORDER_LEFT] = &Polygon{
        Id: WINDOW_BORDER_LEFT,
        Origin: Vector{ x: 0, y: 0 },
        Vertices: []Vector{
            Vector{ x: 0, y: 0 },
            Vector{ x: 0, y: float64(clientRect.Bottom) },
            Vector{ x: -1, y: float64(clientRect.Bottom) },
            Vector{ x: -1, y: 0 },
        },
        Renderer: renderers[RENDER_BG],
    }
    objects[WINDOW_BORDER_BOTTOM] = &Polygon{
        Id: WINDOW_BORDER_BOTTOM,
        Origin: Vector{ x: 0, y: 0 },
        Vertices: []Vector{
            Vector{ x: 0, y: float64(clientRect.Bottom) },
            Vector{ x: float64(clientRect.Right), y: float64(clientRect.Bottom) },
            Vector{ x: float64(clientRect.Right), y: float64(clientRect.Bottom) + 1},
            Vector{ x: 0, y: float64(clientRect.Bottom) + 1},
        },
        Renderer: renderers[RENDER_BG],
    }
    objects[WINDOW_BORDER_RIGHT] = &Polygon{
        Id: WINDOW_BORDER_RIGHT,
        Origin: Vector{ x: 0, y: 0 },
        Vertices: []Vector{
            Vector{ x: float64(clientRect.Right), y: 0 },
            Vector{ x: float64(clientRect.Right) + 1, y: 0 },
            Vector{ x: float64(clientRect.Right) + 1, y: float64(clientRect.Bottom) },
            Vector{ x: float64(clientRect.Right), y: float64(clientRect.Bottom) },
        },
        Renderer: renderers[RENDER_BG],
    }
}

func PhysicsTick(hwnd w32.HWND) {
    UpdateSpeedsAndPositions()
    UpdateCollisions()
}
