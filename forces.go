package gologo

import (
    //"fmt"
    "math"
)

type ForceGenerator interface {
    UpdateForce(*Particle)
}

var gravity = &GravityGenerator {
    Acceleration: Vector2 { 0, 0 },
}

type GravityGenerator struct {
    Acceleration Vector2
}

func (g *GravityGenerator) UpdateForce(p *Particle) {
    if (math.IsInf(p.GetMass(), 1)) {
        return
    }

    force := g.Acceleration
    force.Scale(p.GetMass())
    p.AddForce(&force)
} 

type AttractorGenerator struct {
    Position Vector2
    Factor float64
}

func (g *AttractorGenerator) UpdateForce(p *Particle) {
    if (math.IsInf(p.GetMass(), 1)) {
        return
    }

    direction := DirectionVector(&p.Position, &g.Position)
    //fmt.Printf("calculating attractor force: direction vector: %v\n", direction)
    force := UnitVector(direction)
    //fmt.Printf("calculating attractor force: unit vector: %v\n", force)
    //fmt.Printf("calculating attractor force: scaling factor: %v - factor: %v, mass: %v, square magnitude: %v\n", g.Factor * obj.GetMass() * SquareMagnitude(direction), g.Factor, obj.GetMass(), SquareMagnitude(direction))
    force.Scale(g.Factor * p.GetMass() * Magnitude(direction))
    //fmt.Printf("Calculating force: force: %v, dist ^ 2: %v\n", force, Magnitude(direction))
    p.AddForce(force)
}

type DragGenerator struct {
    k1 float64
    k2 float64
}

func (g *DragGenerator) UpdateForce(p *Particle) {
    if (math.IsInf(p.GetMass(), 1)) {
        return
    }

    force := &p.Velocity

    if (force.x > 0 || force.y > 0) {        
        magnitude := Magnitude(force)
        drag := g.k1 * magnitude + g.k2 * magnitude * magnitude

        force = UnitVector(force)
        force.Scale(-drag)
        p.AddForce(force)
    }
}