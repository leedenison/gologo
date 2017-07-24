package gologo

import (
    "fmt"
    "github.com/go-gl/mathgl/mgl32"
)

type ContactGenerator interface {
    GenerateContacts() []*Contact
    UpdateContact(contact *Contact, contactIdx int, resolved *Contact, resolvedIdx int)
}

/////////////////////////////////////////////////////////////
// TaggedContactGenerator
//

type TaggedContactGenerator struct {
    SourceTag string
    TargetTag string
    PenetrationResolver PenetrationResolver
    PostContactResolver PostContactResolver
}

func (cg *TaggedContactGenerator) GenerateContacts() []*Contact {
    result := []*Contact {}
    sourceSet, existsSrc := Tags[cg.SourceTag]
    targetSet, existsTgt := Tags[cg.TargetTag]

    if existsSrc && existsTgt {
        for objectA, _ := range sourceSet {
            for objectB, _ := range targetSet {
                if objectA != objectB {
                    contactPoint, contactNormal, penetration :=
                        cg.GenerateContactData(objectA, objectB)

                    if penetration > 0 {
                        result = append(result, &Contact {
                            Objects: [2]*Object { objectA, objectB },
                            ContactGenerator: cg,
                            ContactNormal: contactNormal,
                            ContactPoint: contactPoint,
                            Penetration: penetration,
                            PenetrationResolver: cg.PenetrationResolver,
                            PostContactResolver: cg.PostContactResolver,
                        })
                    }
                }
            }
        }
    }

    return result
}

func (cg *TaggedContactGenerator) UpdateContact(
        contact *Contact, contactIdx int, resolved *Contact, resolvedIdx int) {
    if contact.Objects[0] == nil || contact.Objects[1] == nil {
        panic(fmt.Sprintf("Contact must involve two objects: %v\n", contact))
    }
    contact.ContactPoint, contact.ContactNormal, contact.Penetration =
        cg.GenerateContactData(contact.Objects[0], contact.Objects[1])
}

func (cg *TaggedContactGenerator) GenerateContactData(
        objectA *Object, objectB *Object) (mgl32.Vec4, mgl32.Vec4, float32) {
    switch primitiveA := objectA.ObjectType.Primitive.(type) {
    case *Circle:
        switch primitiveB := objectB.ObjectType.Primitive.(type) {
        case *Circle:
            return CalcCircleCircleContact(
                primitiveA, objectA.Model, primitiveB, objectB.Model)
        default:
            panic(fmt.Sprintf("Unhandled primitive type: %t\n", objectB.ObjectType.Primitive))
        }
    default:
        panic(fmt.Sprintf("Unhandled primitive type: %t\n", objectA.ObjectType.Primitive))
    }
}

/////////////////////////////////////////////////////////////
// ScreenEdgeContactGenerator
//

type ScreenEdgeContactGenerator struct {
    Tag string
    PenetrationResolver PenetrationResolver
    PostContactResolver PostContactResolver
}

type ScreenEdgeContactData struct {
    Direction ScreenDirection
}

func (cg *ScreenEdgeContactGenerator) GenerateContacts() []*Contact {
    result := []*Contact {}
    set, exists := Tags[cg.Tag]

    if exists {
        for object, _ := range set {
            for direction := SCREEN_UP; direction <= SCREEN_RIGHT; direction++ {
                contactPoint, contactNormal, penetration :=
                    cg.GenerateContactData(object, direction)

                if penetration > 0 {
                    result = append(result, &Contact {
                        Objects: [2]*Object { object },
                        ContactGenerator: cg,
                        GeneratorData: &ScreenEdgeContactData { direction },
                        ContactNormal: contactNormal,
                        ContactPoint: contactPoint,
                        Penetration: penetration,
                        PenetrationResolver: cg.PenetrationResolver,
                        PostContactResolver: cg.PostContactResolver,
                    })
                }
            }
        }
    }

    return result
}

func (cg *ScreenEdgeContactGenerator) UpdateContact(
        contact *Contact, contactIdx int, resolved *Contact, resolvedIdx int) {
    if contact.Objects[0] == nil || contact.Objects[1] != nil {
        panic(fmt.Sprintf("Contact must involve exactly one object: %v\n", contact))
    }

    generatorData, ok := contact.GeneratorData.(*ScreenEdgeContactData)
    if !ok {
        panic(fmt.Sprintf("Invalid ScreenEdgeContactData: %v\n", generatorData))
    }

    contact.ContactPoint, contact.ContactNormal, contact.Penetration =
        cg.GenerateContactData(contact.Objects[0], generatorData.Direction)
}

func (cg *ScreenEdgeContactGenerator) GenerateContactData(
        object *Object, direction ScreenDirection) (mgl32.Vec4, mgl32.Vec4, float32) {
    switch primitive := object.ObjectType.Primitive.(type) {
    case *Circle:
        position := object.Model.Col(3)
        switch direction {
        case SCREEN_UP:
            return mgl32.Vec4 {
                position.X(),
                (float32(glWin.Height) + position.Y() + primitive.Radius) / 2,
                0.0,
                1.0,
            },
            mgl32.Vec4 { 0.0, -1.0, 0.0, 1.0 },
            position.Y() + primitive.Radius - float32(glWin.Height)
        case SCREEN_DOWN:
            return mgl32.Vec4 {
                position.X(),
                (position.Y() - primitive.Radius) / 2,
                0.0,
                1.0,
            },
            mgl32.Vec4 { 0.0, 1.0, 0.0, 1.0 },
            primitive.Radius - position.Y()
        case SCREEN_LEFT:
            return mgl32.Vec4 {
                (position.X() - primitive.Radius) / 2,
                position.Y(),
                0.0,
                1.0,
            },
            mgl32.Vec4 { 1.0, 0.0, 0.0, 1.0 },
            primitive.Radius - position.X()
        case SCREEN_RIGHT:
            return mgl32.Vec4 {
                (float32(glWin.Width) + position.X() + primitive.Radius) / 2,
                position.Y(),
                0.0,
                1.0,
            },
            mgl32.Vec4 { -1.0, 0.0, 0.0, 1.0 },
            position.X() + primitive.Radius - float32(glWin.Width)
        default:
            panic(fmt.Sprintf("Unknown ScreenDirection: %v\n", direction))
        }
    default:
        panic(fmt.Sprintf("Unhandled primitive type: %t\n", object.ObjectType.Primitive))
    }
}
