package gologo

type ContactGenerator interface {
    AddContacts() []*ParticleContact
}

type ParticleContact struct {
    Particles [2]*Particle
    Restitution float64
    ContactNormal Vector2
    Penetration float64
}

func (p *ParticleContact) CalcSeparatingVelocity() float64 {
    relative := p.Particles[0].Velocity

    if (p.Particles[1] != nil) {
        relative.Subtract(&p.Particles[1].Velocity)
    }

    return DotProduct(&relative, &p.ContactNormal)
}

func (p *ParticleContact) ResolveVelocity() {
    separationV := p.CalcSeparatingVelocity()

    if (separationV > 0) {
        return
    }

    // Reverse the direction and absorb some bounce
    newSeparationV := -separationV * p.Restitution

    // Calculate velocity due to acceleration
    accelerationV := p.Particles[0].CalcAcceleration()
    if (p.Particles[1] != nil) {
        accelerationV.Subtract(p.Particles[1].CalcAcceleration())
    }

    // Calculate separation velocity due to acceleration
    separationAccelerationV := DotProduct(accelerationV, &p.ContactNormal)

    // We decreased separation velocity by the amount due to
    // forces acting on the particles in the last frame, if the 
    // actual separation velocity was less than this amount we
    // we fix it at zero
    if (separationAccelerationV < 0) {
        newSeparationV += p.Restitution * separationAccelerationV

        if (newSeparationV < 0) {
            newSeparationV = 0
        }
    }

    deltaV := newSeparationV - separationV 

    totalInverseMass := p.Particles[0].InverseMass
    if (p.Particles[1] != nil) {
        totalInverseMass += p.Particles[1].InverseMass
    }

    // If all particles involved have infinite mass
    // impulses have no effect
    if (totalInverseMass <= 0) {
        return
    }

    // Impulse to apply
    impulse := deltaV / totalInverseMass

    impulsePerInverseMass := p.ContactNormal
    impulsePerInverseMass.Scale(impulse)

    // Apply impulses in the direction of the contact
    p.Particles[0].Velocity.AddScaledVector(
        &impulsePerInverseMass,
        p.Particles[0].InverseMass)

    if (p.Particles[1] != nil) {
        // Apply negatively scaled impulse to reverse
        // the direction
        p.Particles[1].Velocity.AddScaledVector(
            &impulsePerInverseMass,
            -p.Particles[1].InverseMass)
    }
}

func (p *ParticleContact) ResolveInterpenetration() {
    if (p.Penetration <= 0) {
        return
    }

    totalInverseMass := p.Particles[0].InverseMass

    if (p.Particles[1] != nil) {
        totalInverseMass += p.Particles[1].InverseMass
    }

    if (totalInverseMass <= 0) {
        return
    }

    movePerInverseMass := p.ContactNormal
    movePerInverseMass.Scale(p.Penetration / totalInverseMass)

    movement := movePerInverseMass
    movement.Scale(p.Particles[0].InverseMass) 
    p.Particles[0].Position.Add(&movement)

    if (p.Particles[1] != nil) {
        movement = movePerInverseMass
        movement.Scale(-p.Particles[1].InverseMass)
        p.Particles[1].Position.Add(&movement)        
    }
}

type ParticleRod struct {
    Particles [2]*Particle
    Length float64
}

func (p *ParticleRod) CurrentLength() float64 {
    relativePos := p.Particles[0].Position
    relativePos.Subtract(&p.Particles[1].Position)
    return Magnitude(&relativePos)
}

func (p *ParticleRod) AddContacts() []*ParticleContact {
    result := []*ParticleContact {}
    length := p.CurrentLength()

    if (length == p.Length) {
        return result
    }

    relativePos := p.Particles[1].Position
    relativePos.Subtract(&p.Particles[0].Position)
    normal := UnitVector(&relativePos)
    penetration := length - p.Length

    if (length <= p.Length) {
        normal.Scale(-1)
        penetration = -penetration
    }

    result = append(result, &ParticleContact{
        Particles: p.Particles,
        ContactNormal: *normal,
        Penetration: penetration,
        Restitution: 0,
    })

    return result
}

type ParticleCable struct {
    Particles [2]*Particle
    MaxLength float64
    Restitution float64
}

func (p *ParticleCable) CurrentLength() float64 {
    relativePos := p.Particles[0].Position
    relativePos.Subtract(&p.Particles[1].Position)
    return Magnitude(&relativePos)
}

func (p *ParticleCable) AddContacts() []*ParticleContact {
    result := []*ParticleContact {}
    length := p.CurrentLength()

    if (length < p.MaxLength) {
        return result
    }

    relativePos := p.Particles[1].Position
    relativePos.Subtract(&p.Particles[0].Position)
    normal := UnitVector(&relativePos)

    result = append(result, &ParticleContact{
        Particles: p.Particles,
        ContactNormal: *normal,
        Penetration: length - p.MaxLength,
        Restitution: p.Restitution,
    })

    return result
}