package gologo

import (
    "fmt"
    "math"
    "github.com/go-gl/mathgl/mgl32"
    "strings"
)

var objectIndex = map[*Object]*Thing {}

func OnHit(tag1 string, tag2 string, f func(*Thing, *Thing)) {
    post := &ThingPostContactResolver {
        Callback: f,
    }

    CreateTaggedContactGenerator(tag1, tag2, nil, post)
}

type ThingPostContactResolver struct {
    Callback func(*Thing, *Thing)
}

func (t *ThingPostContactResolver) ResolveContact(contact *Contact) {
    thing1, thing1Exists := objectIndex[contact.Objects[0]]
    if !thing1Exists {
        thing1 = &Thing {}
    }

    thing2, thing2Exists := objectIndex[contact.Objects[1]]
    if !thing2Exists {
        thing2 = &Thing {}
    }

    t.Callback(thing1, thing2)
}

func ShowAllThings(prefix string) {
    objectSpace := 100
    objectsPerRow := (DEFAULT_WIN_SIZE_X / objectSpace) - 1
    i := 0

    for name, _ := range configs {
        if strings.HasPrefix(name, prefix) {
            positionX := (i % objectsPerRow + 1) * objectSpace
            positionY := (i / objectsPerRow + 1) * objectSpace

            Trace.Printf("Building thing(%v) at: %v, %v\n", name, positionX, positionY)
            Builder().SetPosition(positionX, positionY).Build(name)

            i++
        }
    }
}

/////////////////////////////////////////////////////////////
// Things
//

type Thing struct {
    Object *Object
}

func (t *Thing) GetAge() int {
    if t.Object == nil {
        return 0
    }
    return GetTickTime() - t.Object.Creation
}

func (t *Thing) GetPosition() (int, int) {
    if t.Object == nil {
        return 0, 0
    }
    return int(t.Object.Model.Col(3).X()), int(t.Object.Model.Col(3).Y())
}

func (t *Thing) SetPosition(x int, y int) {
    if t.Object == nil {
        return
    }
    t.Object.Model.SetCol(3, mgl32.Vec4 { float32(x), float32(y), 0.0, 1.0 })
}

func (t *Thing) SetPositionVec2(p mgl32.Vec2) {
    if t.Object == nil {
        return
    }
    t.Object.Model.SetCol(3, p.Vec4(0.0, 1.0))
}

func (t *Thing) SetZOrder(z int) {
    if t.Object == nil {
        return
    }
    t.Object.ZOrder = z
}

func (t *Thing) MoveForward(amount int) {
    if t.Object == nil {
        return
    }
    forward := t.Object.Model.Col(1).Vec3().Normalize()
    t.Object.Model = mgl32.Translate3D(
            forward.X() * float32(amount),
            forward.Y() * float32(amount),
            0.0).
        Mul4(t.Object.Model)
}


func (t *Thing) MoveBack(amount int) {
    if t.Object == nil {
        return
    }
    forward := t.Object.Model.Col(1).Vec3().Normalize()
    t.Object.Model = mgl32.Translate3D(
            -forward.X() * float32(amount),
            -forward.Y() * float32(amount),
            0.0).
        Mul4(t.Object.Model)
}


func (t *Thing) MoveLeft(amount int) {
    if t.Object == nil {
        return
    }
    right := t.Object.Model.Col(0).Vec3().Normalize()
    t.Object.Model = mgl32.Translate3D(
            -right.X() * float32(amount),
            -right.Y() * float32(amount),
            0.0).
        Mul4(t.Object.Model)
}


func (t *Thing) MoveRight(amount int) {
    if t.Object == nil {
        return
    }
    right := t.Object.Model.Col(0).Vec3().Normalize()
    t.Object.Model = mgl32.Translate3D(
            right.X() * float32(amount),
            right.Y() * float32(amount),
            0.0).
        Mul4(t.Object.Model)
}

func (t *Thing) TurnClockwise(angle int) {
    rotation := mgl32.HomogRotate3DZ(
        mgl32.DegToRad(float32(angle)))

    t.Object.Model = t.Object.Model.Mul4(rotation)
}

func (t *Thing) TurnAntiClockwise(angle int) {
    rotation := mgl32.HomogRotate3DZ(
        mgl32.DegToRad(float32(-angle)))

    t.Object.Model = t.Object.Model.Mul4(rotation)
}

func (t *Thing) Direction() int {
    angle := math.Atan2(float64(t.Object.Model.At(1, 1)), float64(t.Object.Model.At(0, 1))) - math.Pi / 2

    if angle <= -0 {
        return int(mgl32.RadToDeg(float32(angle + 2 * math.Pi)))
    } else {
        return int(mgl32.RadToDeg(float32(angle)))
    }
}

func (t *Thing) DirectionOf(other *Thing) int {
    direction := other.Object.Model.Col(3).Vec3().Sub(t.Object.Model.Col(3).Vec3())
    angle := math.Atan2(float64(direction[1]), float64(direction[0])) - math.Pi / 2

    if angle <= -0 {
        return int(mgl32.RadToDeg(float32(angle + 2 * math.Pi)))
    } else {
        return int(mgl32.RadToDeg(float32(angle)))
    }
}

func (t *Thing) IsOnScreen() bool {
    if t.Object == nil {
        return false
    }
    switch primitive := t.Object.Primitive.(type) {
    case nil:
        return OriginIsOnScreen(t.Object.Model)
    case *Circle:
        return CircleIsOnScreen(primitive, t.Object.Model)
    default:
        panic(fmt.Sprintf("Unhandled primitive type: %t\n", t.Object.Primitive))
    }
}

func (t *Thing) Delete() {
    if t.Object == nil {
        return
    }
    UntagAll(t.Object)

    for idx, object := range objects {
        if object == t.Object {
            if len(objects) > 1 {
                objects = append(objects[:idx], objects[idx+1:]...)
            } else {
                objects = objects[0:0]
            }
        }
    }

    delete(objectIndex, t.Object)
    t.Object = nil
}

func (t *Thing) IsDeleted() bool {
    return t.Object == nil
}

/////////////////////////////////////////////////////////////
// ThingBuilder
//

type ThingBuilder struct {
    Config string
    Position mgl32.Mat4
    Orientation mgl32.Mat4
    RenderScale mgl32.Mat4
    ZOrder int
    Tags []string
    RenderData interface {}
}

func Builder() *ThingBuilder {
    return &ThingBuilder {
        Position: DEFAULT_POSITION,
        Orientation: DEFAULT_ORIENTATION,
        RenderScale: DEFAULT_SCALE,
    }
}

func (sb *ThingBuilder) SetDepth(z float32) *ThingBuilder {
    position := sb.Position.Col(3)
    position[2] = z
    sb.Position.SetCol(3, position)
    return sb
}

func (sb *ThingBuilder) SetScale(factor float32) *ThingBuilder {
    sb.RenderScale = mgl32.Scale3D(factor, factor, 1)
    return sb
}

func (sb *ThingBuilder) SetZOrder(z int) *ThingBuilder {
    sb.ZOrder = z
    return sb
}

func (sb *ThingBuilder) SetDirection(angle int) *ThingBuilder {
    sb.Orientation = mgl32.HomogRotate3DZ(
        mgl32.DegToRad(float32(angle)))
    return sb
}

func (sb *ThingBuilder) SetPosition(x, y int) *ThingBuilder {
    sb.Position = mgl32.Translate3D(float32(x), float32(y), 0.0)
    return sb
}

func (sb *ThingBuilder) AddTag(tag string) *ThingBuilder {
    sb.Tags = append(sb.Tags, tag)
    return sb
}

func (sb *ThingBuilder) Build(thingType string) *Thing {
    model := sb.Position.Mul4(sb.Orientation.Mul4(sb.RenderScale))

    template, ok := templates[thingType]
    if !ok {
        panic(fmt.Sprintf("Invalid object template: %v\n", thingType))
    }

    object := template.CreateObject(model)
    object.ZOrder = sb.ZOrder

    meshRenderer, ok := object.Renderer.(*MeshRenderer)
    if object.Primitive == nil && ok {
        object.Primitive = InitCircleFromMesh(meshRenderer.MeshVertices)
    }

    objects = append(objects, object)

    for _, tag := range sb.Tags {
        Tag(object, tag)
    }

    objectIndex[object] = &Thing {
        Object: object,
    }
    return objectIndex[object]
 }

/////////////////////////////////////////////////////////////
// ThingList
//

type ThingList struct {
    Data []*Thing
}

func (t *ThingList) RemoveAt(i int) {
    if len(t.Data) > 1 {
        t.Data = append(t.Data[:i], t.Data[i+1:]...)
    } else {
        t.Data = t.Data[0:0]
    }
}

func (t *ThingList) Remove(thing *Thing) {
    for i := 0; i < t.Length(); i++ {
        listThing := t.Get(i)
        if listThing == thing {
            t.RemoveAt(i)
            i--
        }
    }
}

func (t *ThingList) Add(thing *Thing) {
    t.Data = append(t.Data, thing)
}

func (t *ThingList) Length() int {
    return len(t.Data)
}

func (t *ThingList) Get(i int) *Thing {
    return t.Data[i]
}

func (t *ThingList) Contains(thing *Thing) bool {
    for _, listThing := range t.Data {
        if thing == listThing {
            return true
        }
    }

    return false
}

/////////////////////////////////////////////////////////////
// TextBuilder
//

func (sb *ThingBuilder) BuildText(text string) *Thing {
    thing := sb.Build("TEXT");
    renderer := thing.Object.Renderer

    textRenderer, ok := renderer.(*TextRenderer)
    if !ok {
        panic(fmt.Sprintf("Invalid text renderer type: %t\n", renderer))
    }

    textRenderer.Text = []byte(text)
    return thing
}
