package gologo

import (
    //"fmt"
    "math"
)

type Object interface {
    GetId() int
    GetRenderer() Renderer
    Integrate()
}

type Particle struct {
    Position Vector2
    Velocity Vector2
    LastAcceleration Vector2
    Force Vector2
    Damping float64
    InverseMass float64
}

type IParticle interface {
    GetParticle() *Particle
}

func (p *Particle) AddForce(force *Vector2) {
    p.Force.Add(force)
}

func (p *Particle) ClearForces() {
    p.Force = Vector2 { 0, 0 }
}

func (p *Particle) CalcForce() *Vector2 {
    if generators, ok := forceRegistry[p]; ok {
        for _, generator := range generators {
            generator.UpdateForce(p)
        }
    }
    return &p.Force
}

func (p *Particle) CalcAcceleration() *Vector2 {
    acceleration := p.Force
    acceleration.Scale(p.InverseMass)
    return &acceleration
}

func (p *Particle) SetLastAcceleration(a *Vector2) {
    p.LastAcceleration = *a
}

func (p *Particle) GetMass() float64 {
    if (p.InverseMass <= 0.0) {
        return math.Inf(1)
    }

    return 1 / p.InverseMass
}

// Integrate using the velocity verlet method
func (p *Particle) Integrate() {
    if (p.InverseMass <= 0.0) {
        return
    }

    // oldPosition := c.Position
    // oldAcceleration := c.LastAcceleration
    // oldVelocity := c.Velocity

    // We update position based on velocity at the
    // beginning of the tick and the average
    // acceleration from the last tick:
    //   x += ut + 0.5 * a * t ^ 2
    p.Position.Add(&p.Velocity)
    p.Position.AddScaledVector(&p.LastAcceleration, 0.5)
    // fmt.Printf("Integrate: Update position (v + 0.5 * a)\t{ Id: %v, Position: %.2f, Velocity: %.2f, Force: %.2f, Acceleration: %.2f, LastAcceleration: %.2f }\n",
    //     c.Id, c.Position, c.Velocity, c.Force, c.GetAcceleration(), c.LastAcceleration)
    // newPosition := c.Position
    // newPosition.Subtract(&oldPosition)
    // fmt.Printf("Integrate: Update\t{ dP: %.3f, ", newPosition)

    // Calculate the average acceleration for this tick
    p.CalcForce()
    p.LastAcceleration.Add(p.CalcAcceleration())
    p.LastAcceleration.Scale(0.5)
    p.ClearForces()

    // fmt.Printf("Integrate: Update acceleration ((a + a') * 0.5)\t{ Id: %v, Position: %.2f, Velocity: %.2f, Force: %.2f, Acceleration: %.2f, LastAcceleration: %.2f }\n",
    //     c.Id, c.Position, c.Velocity, c.Force, c.GetAcceleration(), c.LastAcceleration)
    // newAcceleration := c.LastAcceleration
    // newAcceleration.Subtract(&oldAcceleration)
    // fmt.Printf("dA: %.3f, ", newAcceleration)

    // Update velocity based on average acceleration.
    //   v += 0.5 * at
    p.Velocity.Add(&p.LastAcceleration)
    p.Velocity.Scale(p.Damping)

    // fmt.Printf("Integrate: Update velocity (v + (a + a') * 0.5)\t{ Id: %v, Position: %.2f, Velocity: %.2f, Force: %.2f, Acceleration: %.2f, LastAcceleration: %.2f }\n",
    //     c.Id, c.Position, c.Velocity, c.Force, c.GetAcceleration(), c.LastAcceleration)
    // newVelocity := c.Velocity
    // newVelocity.Subtract(&oldVelocity)
    // fmt.Printf("dV: %.3f }\n", newVelocity)
}

type Circle struct {
    Id int
    Renderer Renderer
    Radius float64
    Particle Particle
}

func (c *Circle) GetId() int {
    return c.Id
}

func (c *Circle) GetRenderer() Renderer {
    return c.Renderer
}

func (c *Circle) Integrate() {
    c.Particle.Integrate()
}

func (c *Circle) GetParticle() *Particle {
    return &c.Particle
}