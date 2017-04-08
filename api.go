package gologo

const DEFAULT_RADIUS = 20

func CreateCircle(x, y int32) *Circle {
    return CreateCircleWithSpeed(x, y, 0, 0)
}

func CreateCircleWithSpeed(x, y, speedx, speedy int32) *Circle {
    obj := &Circle{
        Id: nextObjectId,
        Renderer: &CIRCLE_RENDERER,
        Radius: DEFAULT_RADIUS,
        Particle: Particle{
            Position: Vector2{ x: float64(x), y: float64(y) },
            Velocity: Vector2{ x: float64(speedx) / SPEED_FACTOR, y: float64(speedy) / SPEED_FACTOR },
            Force: Vector2{ x: 0.0, y: 0.0 },
            InverseMass: AREA_TO_MASS_RATIO / CircleArea(DEFAULT_RADIUS),
            Damping: DEFAULT_DAMPING,
        },
    }
    objects[nextObjectId] = obj
    nextObjectId++
    return obj
}

func AttachCable(p1 IParticle, p2 IParticle, length int32) {
    contactRegistry = append(
        contactRegistry,
        &ParticleCable {
            Particles: [2]*Particle {
                p1.GetParticle(),
                p2.GetParticle(),
            },
            MaxLength: float64(length),
            Restitution: 0.2,
        })
}

func AttachFixedSpring(x, y int32, p IParticle) {
    forceRegistry[p.GetParticle()] = append(
        forceRegistry[p.GetParticle()],
        &AttractorGenerator {
            Position: Vector2 { float64(x), float64(y) },
            Factor: 0.001,
        })
}

func AttachGravity(p IParticle) {
    forceRegistry[p.GetParticle()] = append(forceRegistry[p.GetParticle()], gravity)
}

func AttachDrag(p IParticle) {
    forceRegistry[p.GetParticle()] = append(
        forceRegistry[p.GetParticle()],
        &DragGenerator {
            k1: 0.0005,
            k2: 0.0015,
        })
}

func SetGravity(x, y int32) {
    gravity.Acceleration.x = float64(x) / SPEED_FACTOR
    gravity.Acceleration.y = float64(y) / SPEED_FACTOR
}
