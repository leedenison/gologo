package gologo

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestParticleContact_CalcSeparatingVelocity_distContactStationary(t *testing.T) {
    assert := assert.New(t)

    particle1 := particle_m10_p0_0_z0
    particle2 := particle_m10_p0_10_z0

    contact := ParticleContact {
        Particles: [2]*Particle {
            &particle1,
            &particle2,
        },
    }
    separatingV := contact.CalcSeparatingVelocity()

    assert.Equal(float64(0), separatingV, "Separating velocity of stationary objects should be 0.")
}

func TestParticleContact_CalcSeparatingVelocity_movingDirectlyApartWithInverseContactNormal(t *testing.T) {
    assert := assert.New(t)

    particle1 := particle_m10_p0_0_z0
    particle2 := particle_m10_p0_11_z0
    particle2.Velocity = Vector2 { 0, 1 }

    contact := ParticleContact {
        Particles: [2]*Particle {
            &particle1,
            &particle2,
        },
        Restitution: 0.5,
        ContactNormal: Vector2 { 0, 1 },
        Penetration: 1.0,
    }
    separatingV := contact.CalcSeparatingVelocity()

    assert.Equal(float64(-1.0), separatingV, "Separating velocity should be -1.0.")
}

func TestParticleContact_CalcSeparatingVelocity_movingObliquelyTowards(t *testing.T) {
    assert := assert.New(t)

    particle1 := particle_m10_p0_0_z0
    particle1.Velocity = Vector2 { 1, 1 }
    particle2 := particle_m10_p0_10_z0
    particle2.Velocity = Vector2 { 1, -1 }

    contact := ParticleContact {
        Particles: [2]*Particle {
            &particle1,
            &particle2,
        },
        Restitution: 0.5,
        ContactNormal: Vector2 { 0, -1 },
        Penetration: 1.0,
    }
    separatingV := contact.CalcSeparatingVelocity()

    assert.Equal(float64(-2.0), separatingV, "Separating velocity should be -2.0.")
}

func TestParticleContact_ResolveVelocity_restingContact(t *testing.T) {
    assert := assert.New(t)

    particle1 := particle_m10_p0_0_z0

    particle2 := particle_m10_p0_10_z0
    particle2.Velocity = Vector2 { 0, -0.5 }
    particle2.Force = Vector2 {
        x: 0,
        y: -1.0 / particle2.InverseMass,
    }

    contact := ParticleContact {
        Particles: [2]*Particle {
            &particle1,
            &particle2,
        },
        Restitution: 0.5,
        ContactNormal: Vector2 { 0, -1 },
        Penetration: 0.5,
    }
    
    contact.ResolveVelocity()

    assert.Equal(Vector2 { 0, -0.25 }, particle1.Velocity, "particle1 velocity should be { 0, -0.25}.")
    assert.Equal(Vector2 { 0, -0.25 }, particle2.Velocity, "particle2 velocity should be { 0, -0.25}.")
}

func TestParticleContact_ResolveInterpenetration_restingContact(t *testing.T) {
    assert := assert.New(t)

    particle1 := particle_m10_p0_0_z0

    particle2 := particle_m10_p0_10_z0
    particle2.Velocity = Vector2 { 0, -0.5 }
    particle2.Force = Vector2 {
        x: 0,
        y: -1.0 / particle2.InverseMass,
    }

    contact := ParticleContact {
        Particles: [2]*Particle {
            &particle1,
            &particle2,
        },
        Restitution: 0.5,
        ContactNormal: Vector2 { 0, -1 },
        Penetration: 0.5,
    }
    
    contact.ResolveInterpenetration()

    assert.Equal(Vector2 { 0, -0.25 }, particle1.Position, "particle1 position should be { 0, -0.25}.")
    assert.Equal(Vector2 { 0, 10.25 }, particle2.Position, "particle2 position should be { 0, 10.25}.")
}
 