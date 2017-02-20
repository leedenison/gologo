package gologo

import (
    "C"
    "math"
    "github.com/AllenDang/w32"
)

const PHYSICS_MULT = 2
const SPEED_MULT = PHYSICS_MULT * 10

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

var nextObjectId = 4
var objects = map[int]Object {}
var movables = map[int]Vector {}

var gravity = Vector { 0, 0 }

type Vector struct {
    x, y   float64
}

type Object interface {
    GetOrigin() Vector
    SetOrigin(Vector)
}

type Collision struct {
    collisionType int
    impactObj1 Vector
    impactObj2 Vector
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

// TODO: Find the closest collision instead of the first, ie. we could be overlapping two 
//       edges or a corner.  We should calculate the nearer of the edges to bounce off of.
// 
// TODO: Collisions should return three impact points - the point of impact on each shape
//       and the point in world coordinates of the impact.  The intention is that the 
//       points of impact on each shape will be translated to the world co-ordinates point
//       and the rebound is calculated assuming that was the location of the impact.
func CheckCollisionPolyCirc(poly *Polygon, circ *Circle) *Collision {
    //fmt.Printf("Checking collision between:\n\tPolygon: %v\n\tCircle: %v\n", *poly, *circ)
    isRight := true
    vertices := poly.GetWorldVertices()

    for idx, v1 := range vertices {
        idx2 := (idx + 1) % len(vertices)
        v2 := vertices[idx2]

        // Calculate direction vectors for v1 -> c and v1 -> v2, then
        // calculate the dot product
        v3 := DirectionVector(v1, circ.Centre)
        v4 := DirectionVector(v1, v2)
        dot1 := DotProduct(v3, v4)

        // Calculate the square magnitude of the radius and
        // the v1 -> v2 direction vector
        radiusSq := math.Pow(circ.Radius, 2)
        magSq := math.Pow(v4.x, 2) + math.Pow(v4.y, 2)

        // Calculate the distance of the projection of the centre of the
        // circle onto the v1 -> v2 direction vector in units of the 
        // magnitude of the v1 -> v2 direction vector.
        // Clamp the result to limit to the length of the v1 -> v2 direction
        // vector
        projDist := Clamp(dot1 / magSq, 0, 1)

        // Calculate the position of the point on the v1 -> v2 direction
        // vector closes to the projection of the centre of the circle.
        projection := Vector {
            v1.x + projDist * v4.x,
            v1.y + projDist * v4.y,
        }

        v5 := DirectionVector(circ.Centre, projection)
        magSqProj := math.Pow(v5.x, 2) + math.Pow(v5.y, 2)
    
        if magSqProj < radiusSq {
            v6 := UnitVector(v5)
            return &Collision{
                collisionType: OVERLAP_PARTIAL,
                impactObj1: projection,
                impactObj2: Vector { 
                    x: circ.Centre.x + v6.x * circ.Radius,
                    y: circ.Centre.y + v6.y * circ.Radius,
                },
            }                        
        }

        // Not had collision so record if it's on our right
        if (v4.x * v3.y - v4.y * v3.x) < 0 {
            isRight = false
        }
    }

    if isRight == true {
        return &Collision{
            collisionType: OVERLAP_FULL,
            impactObj1: Vector{ x: 0, y: 0 },
            impactObj2: Vector{ x: 0, y: 0 },
        }
    } else {
        return nil
    }
}

func HandleCollisionPolyCirc(poly *Polygon, circ *Circle, collision *Collision) {
    // Current approach assumes the polygon has infinite mass (ie. cannot be moved)
    // and is stationary (ie. does not add energy to the circle).

    if velocity, ok := movables[circ.Id]; ok {
        initialVelocity := Vector {
            x: velocity.x - gravity.x,
            y: velocity.y - gravity.y,
        }
        
        initialPosition := Vector {
            x: circ.Centre.x - (initialVelocity.x + (gravity.x / 2)),
            y: circ.Centre.y - (initialVelocity.y + (gravity.y / 2)),
        }

        // First translate the circle to the point of impact
        translate := DirectionVector(collision.impactObj2, collision.impactObj1)
        circ.Centre.x += translate.x
        circ.Centre.y += translate.y

        // Calculate unit direction vector impact -> centre
        dirVector := UnitDirectionVector(collision.impactObj1, circ.Centre)

        translationDistance := DirectionVector(initialPosition, circ.Centre)

        impactVelocity := Vector {
            x: math.Sqrt(math.Pow(initialVelocity.x, 2) + 2.0 * gravity.x * translationDistance.x),
            y: math.Sqrt(math.Pow(initialVelocity.y, 2) + 2.0 * gravity.y * translationDistance.y),
        }

        // We assume that the average velocity during the time tick approximates
        // the direction we should reflect.
        if (velocity.x < 0) {
            impactVelocity.x = -impactVelocity.x
        }

        if (velocity.y < 0) {
            impactVelocity.y = -impactVelocity.y
        }

        movables[circ.Id] = ReflectVelocity(dirVector, impactVelocity)
    }
}

// TODO: Calculate if partial or full overlap and closest impact
func CheckCollisionPolyPoly(polyA *Polygon, polyB *Polygon) *Collision {
    polyAVertices := polyA.GetWorldVertices()
    polyBVertices := polyB.GetWorldVertices()

    for _, vertices := range [][]Vector{ polyAVertices, polyBVertices } {
        for idx, v1 := range vertices {
            idx2 := (idx + 1) % len(vertices)
            v2 := vertices[idx2]
            
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
        impactObj1: Vector{ x: 0, y: 0 },
        impactObj2: Vector{ x: 0, y: 0 },
    }
}

func CheckCollisionCircCirc(c1 *Circle, c2 *Circle) *Collision {
    dirVector := DirectionVector(c1.Centre, c2.Centre)
    distSq := math.Pow(dirVector.x, 2) + math.Pow(dirVector.y, 2)
    radiusSq := math.Pow(c1.Radius + c2.Radius, 2)
    collisionType := OVERLAP_PARTIAL

    if distSq <= radiusSq {
        unitVectorA := UnitVector(dirVector)
        unitVectorB := Vector {
            x: -unitVectorA.x,
            y: -unitVectorA.y,
        }

        halfRadiusASq := math.Pow(c1.Radius / 2, 2)
        halfRadiusBSq := math.Pow(c2.Radius / 2, 2)
        if distSq < halfRadiusASq ||
           distSq < halfRadiusBSq {
            collisionType = OVERLAP_FULL
        }

        return &Collision{
            collisionType: collisionType,
            impactObj1: Vector {
                x: unitVectorA.x * c1.Radius,
                y: unitVectorA.y * c1.Radius,
            },
            impactObj2: Vector {
                x: unitVectorB.x * c2.Radius,
                y: unitVectorB.y * c2.Radius,
            },
        }

    }
    return nil
}

func HandleCollisionCircCirc(c1 *Circle, c2 *Circle, collision *Collision) {
    // First translate circle to the point of impact
    translate := DirectionVector(collision.impactObj1, collision.impactObj2)
    c1.Centre.x += translate.x
    c1.Centre.y += translate.y

    // Calculate unit direction vector impact -> centre
    dir1 := UnitDirectionVector(collision.impactObj1, c1.Centre)
    movables[c1.Id] = ReflectVelocity(dir1, movables[c1.Id])

    c2.Centre.x -= translate.x
    c2.Centre.y -= translate.y

    // Reuse dir1 since the result is the same for the reverse direction vector
    movables[c2.Id] = ReflectVelocity(dir1, movables[c2.Id])
}

func UpdateSpeedsAndPositions() {
    for objIdx, speed := range movables {
        obj := objects[objIdx]

        // We update position based on the average
        // speed over the tick
        origin := obj.GetOrigin()
        origin.x += speed.x + (gravity.x / 2)
        origin.y += speed.y + (gravity.y / 2)
        obj.SetOrigin(origin)

        speed.x += gravity.x
        speed.y += gravity.y
        movables[objIdx] = speed
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
                case *Circle:
                    collision := CheckCollisionCircCirc(t1, t2)
                    if collision != nil {
                        HandleCollisionCircCirc(t1, t2, collision)
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
