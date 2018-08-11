package gologo

import (
	"fmt"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

var objectIndex = map[*Object]*Thing{}

func OnHit(tag1 string, tag2 string, f func(*Thing, *Thing)) {
	post := &ThingPostContactResolver{
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
		thing1 = &Thing{}
	}

	thing2, thing2Exists := objectIndex[contact.Objects[1]]
	if !thing2Exists {
		thing2 = &Thing{}
	}

	t.Callback(thing1, thing2)
}

// ShowAllThings : Displays every template object on screen side by side
func ShowAllThings(prefix string) {
	objectSpace := 100
	objectsPerRow := (defaultWinSizeX / objectSpace) - 1
	i := 0

	for name := range configs {
		if strings.HasPrefix(name, prefix) {
			positionX := (i%objectsPerRow + 1) * objectSpace
			positionY := (i/objectsPerRow + 1) * objectSpace

			Trace.Printf("Building thing(%v) at: %v, %v\n", name, positionX, positionY)
			Builder().SetPosition(positionX, positionY).Build(name)

			i++
		}
	}
}

/////////////////////////////////////////////////////////////
// Things
//

// Thing : Wrapper around Object with extra convenience functions
// and simplifications
type Thing struct {
	Object *Object
}

// GetAge : Returns age of Thing since creation
func (t *Thing) GetAge() int {
	if t.Object == nil {
		return 0
	}
	return t.Object.GetAge()
}

// GetPosition : Returns X and Y co-ords of Thing's centre in 2D
func (t *Thing) GetPosition() (int, int) {
	if t.Object == nil {
		return 0, 0
	}
	x, y := t.Object.GetPosition()

	return int(x), int(y)
}

// SetPosition : Sets X and Y co-ords of Thing's centre in 2D
func (t *Thing) SetPosition(x int, y int) {
	if t.Object == nil {
		return
	}
	t.Object.SetPosition(float32(x), float32(y))
}

// SetPositionVec2 : Sets X and Y co-ords of Thing's centre
// in 2D using a vector of 2 elements
func (t *Thing) SetPositionVec2(p mgl32.Vec2) {
	if t.Object == nil {
		return
	}
	t.Object.SetPositionVec2(p)
}

// SetZOrder : Sets the height of Thing in 3D space
// as an integer compared with other objects
func (t *Thing) SetZOrder(z int) {
	if t.Object == nil {
		return
	}
	t.Object.SetZOrder(z)
}

// MoveForward : Move the object forward by amount in the direction it's currently facing
func (t *Thing) MoveForward(amount int) {
	if t.Object == nil {
		return
	}
	forward := t.Object.DirectionNormal()
	t.Object.Translate(forward.X()*float32(amount), forward.Y()*float32(amount))
}

// MoveBack : Move the object backward by amount in the opposite of the direction
// it's currently facing
func (t *Thing) MoveBack(amount int) {
	if t.Object == nil {
		return
	}
	forward := t.Object.DirectionNormal()
	t.Object.Translate(-forward.X()*float32(amount), -forward.Y()*float32(amount))
}

// MoveLeft : Move the object left by amount from the direction it's currently facing
func (t *Thing) MoveLeft(amount int) {
	if t.Object == nil {
		return
	}
	right := t.Object.GetModel().Col(0).Vec3().Normalize()
	t.Object.Translate(-right.X()*float32(amount), -right.Y()*float32(amount))
}

// MoveRight : Move the object right by amount from the direction it's currently facing
func (t *Thing) MoveRight(amount int) {
	if t.Object == nil {
		return
	}

	right := t.Object.GetModel().Col(0).Vec3().Normalize()
	t.Object.Translate(right.X()*float32(amount), right.Y()*float32(amount))
}

// TurnClockwise : Rotates Thing clockwise by angle degrees
func (t *Thing) TurnClockwise(angle int) {
	if t.Object == nil {
		return
	}

	t.Object.Rotate(float32(-angle))
}

// TurnAntiClockwise : Rotates Thing anti-clockwise by angle degrees
func (t *Thing) TurnAntiClockwise(angle int) {
	if t.Object == nil {
		return
	}

	t.Object.Rotate(float32(angle))
}

// Direction : Returns the angle the object has been rotated since it was
// created in degrees
func (t *Thing) Direction() int {
	if t.Object == nil {
		return 0
	}
	angle := t.Object.Direction()

	return radToNearestDeg(angle)
}

// DirectionOf : Calculates the direction in degrees to the passed in Thing
// from the receiving Thing
func (t *Thing) DirectionOf(other *Thing) int {
	if t.Object == nil || other.Object == nil {
		return 0
	}
	angle := t.Object.DirectionOf(other.Object)

	return radToNearestDeg(angle)
}

// IsOnScreen : Returns true if Thing is on screen
func (t *Thing) IsOnScreen() bool {
	if t.Object == nil {
		return false
	}

	rect := Rect{{0, 0}, GetWindowSize()}

	if primitive := t.Object.GetPrimitive(); primitive != nil {
		return primitive.OverlapsWithRect(*t.Object, rect)
	}

	// use if the origin is contained as an approximation
	return t.Object.OriginIsContainedInRect(rect)
}

// Delete : Deletes Thing's object and removes it's tags and it from the object list
// Thing continues to exist
func (t *Thing) Delete() {
	if t.Object == nil {
		return
	}
	UntagAll(t.Object)
	UntagRender(t.Object)
	delete(objectIndex, t.Object)
	t.Object = nil
}

// IsDeleted : Returns true if Thing's object is nil
// i.e. if Thing has been "Delete"d
func (t *Thing) IsDeleted() bool {
	return t.Object == nil
}

/////////////////////////////////////////////////////////////
// ThingBuilder
//

type ThingBuilder struct {
	Config      string
	Position    mgl32.Mat4
	Orientation mgl32.Mat4
	RenderScale mgl32.Mat4
	ZOrder      int
	Tags        []string
	RenderData  interface{}
}

func Builder() *ThingBuilder {
	return &ThingBuilder{
		Position:    defaultPosition,
		Orientation: defaultOrientation,
		RenderScale: defaultScale,
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

	object, err := CreateTemplateObject(thingType, model)
	if err != nil {
		panic(fmt.Sprintf("Failed to create object: %v\n", err))
	}
	object.SetZOrder(sb.ZOrder)

	TagRender(object)

	for _, tag := range sb.Tags {
		Tag(object, tag)
	}

	objectIndex[object] = &Thing{
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
	thing := sb.Build("TEXT")
	renderer := thing.Object.GetRenderer()

	textRenderer, ok := renderer.(*TextRenderer)
	if !ok {
		panic(fmt.Sprintf("Invalid text renderer type: %t\n", renderer))
	}

	textRenderer.Text = []byte(text)
	return thing
}
