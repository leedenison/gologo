package gologo

import (
    //"fmt"
    "math"
    "github.com/AllenDang/w32"
)

const PHYSICS_MULT = 2
const SPEED_FACTOR = PHYSICS_MULT * 10
const MAX_ITERATIONS_MULT = 2
const AREA_TO_MASS_RATIO = 1000
const DEFAULT_DAMPING = 0.9995

const (
    OBJECT_SOLID = iota
    OBJECT_OUTLINE
    OBJECT_EMPTY
)

var nextObjectId = 0
var objects = map[int]Object {}
var forceRegistry = map[*Particle][]ForceGenerator {}
var contactRegistry = []ContactGenerator {}

func Integrate() {
    for _, obj := range objects {
        obj.Integrate()
    }
}

func GenerateContacts() []*ParticleContact {
    result := []*ParticleContact {}

    for _, generator := range contactRegistry {
        contacts := generator.AddContacts()
        result = append(result, contacts...)
    }

    return result
}

func ResolveContacts(contacts []*ParticleContact) {
    for i := 0; i < len(contacts) * MAX_ITERATIONS_MULT; i++ {
        maxSeparationV := math.MaxFloat64
        maxIdx := len(contacts)

        // Find the smallest separation velocity
        for j := 0; j < len(contacts); j++ {
            separationV := contacts[j].CalcSeparatingVelocity()

            if (separationV < maxSeparationV &&
                (separationV < 0 || contacts[j].Penetration > 0)) {
                maxSeparationV = separationV
                maxIdx = j
            }
        }

        if (maxIdx == len(contacts)) {
            break
        }

        contacts[maxIdx].ResolveVelocity()
        contacts[maxIdx].ResolveInterpenetration()
    }
}

func PhysicsTick(hwnd w32.HWND, tickCount uint64) {
    //fmt.Printf("Tick count: %v\n", tickCount)
    Integrate()
    ResolveContacts(GenerateContacts())
}
