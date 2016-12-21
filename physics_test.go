package gologo

import "testing"

func TestCheckCollisionPolyCirc_None(t *testing.T) {
	assert := (*Assert)(t)

    collision := CheckCollisionPolyCirc(SQUARE_150_150, CIRCLE_0_0)

    assert.That(collision).IsNil()
}

func TestCheckCollisionPolyCirc_PartialOverlap(t *testing.T) {
	assert := (*Assert)(t)
    
    collision := CheckCollisionPolyCirc(SQUARE_150_150, CIRCLE_150_210)

    assert.That(collision).Equals(
		&Collision {
			collisionType: OVERLAP_PARTIAL,
			closestImpact: Vector { x: 150, y: 180 },
		})
}

func TestCheckCollisionPolyCirc_FullOverlap(t *testing.T) {
	assert := (*Assert)(t)
    
    collision := CheckCollisionPolyCirc(SQUARE_150_150, CIRCLE_150_150)

    assert.That(collision).Equals(
		&Collision {
			collisionType: OVERLAP_FULL,
			closestImpact: Vector { x: 150, y: 150 },
		})
}


var CIRCLE_0_0 = &Circle{
    Centre: Vector{ x: 0, y: 0 },
    Radius: 20,
}

var CIRCLE_150_210 = &Circle{
    Centre: Vector{ x: 150, y: 210 },
    Radius: 20,
}

var CIRCLE_150_150 = &Circle{
    Centre: Vector{ x: 150, y: 150 },
    Radius: 20,
}

var SQUARE_100 = []Vector{
    Vector{ x: -50, y: -50 },
    Vector{ x: 50, y: -50 },
    Vector{ x: 50, y: 50 },
    Vector{ x: -50, y: 50 },
}

var SQUARE_150_150 = &Polygon{
    Origin: Vector{ x: 150, y: 150 },
    Vertices: SQUARE_100,
}

var SQUARE_350_350 = &Polygon{
    Origin: Vector{ x: 350, y: 350 },
    Vertices: SQUARE_100,
}