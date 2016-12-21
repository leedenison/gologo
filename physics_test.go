package gologo

import "testing"

func TestCheckCollisionPolyCirc_None(t *testing.T) {
	assert := (*Assert)(t)

    collision := CheckCollisionPolyCirc(SQUARE_100_200, CIRCLE_0_0)

    assert.That(collision).IsNil()
}

func TestCheckCollisionPolyCirc_PartialOverlap(t *testing.T) {
	assert := (*Assert)(t)
    
    collision := CheckCollisionPolyCirc(SQUARE_100_200, CIRCLE_150_210)

    assert.That(collision).Equals(
		&Collision {
			collisionType: OVERLAP_PARTIAL,
			closestImpact: Vector { x: 150, y: 180 },
		})
}

func TestCheckCollisionPolyCirc_FullOverlap(t *testing.T) {
	assert := (*Assert)(t)
    
    collision := CheckCollisionPolyCirc(SQUARE_100_200, CIRCLE_150_150)

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

var SQUARE_100_200 = &Polygon{
    Origin: Vector{ x: 100, y: 100 },
    Vertices: []Vector{
        Vector{ x: 100, y: 100 },
        Vector{ x: 200, y: 100 },
        Vector{ x: 200, y: 200 },
        Vector{ x: 100, y: 200 },
    },
}

var SQUARE_300_400 = &Polygon{
    Origin: Vector{ x: 300, y: 300 },
    Vertices: []Vector{
        Vector{ x: 300, y: 300 },
        Vector{ x: 400, y: 300 },
        Vector{ x: 400, y: 400 },
        Vector{ x: 300, y: 400 },
    },
}