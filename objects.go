package gologo

import (
    "github.com/go-gl/mathgl/mgl32"
)

type Object struct {
    Model mgl32.Mat4
    ZOrder int
    Creation int
    Primitive Primitive
    Renderer Renderer
}

func CreateObject(model mgl32.Mat4) *Object {
    return &Object {
        Model: model,
        Creation: GetTickTime(),
    }
}

func (object *Object) Clone() *Object {
    objectCopy := *object

    if object.Renderer != nil {
        objectCopy.Renderer = object.Renderer.Clone()
    }

    if object.Primitive != nil {
        objectCopy.Primitive = object.Primitive.Clone()
    }

    return &objectCopy
}

type ByZOrder []*Object

func (s ByZOrder) Len() int {
    return len(s)
}

func (s ByZOrder) Swap(i int, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s ByZOrder) Less(i int, j int) bool {
    return s[i].ZOrder < s[j].ZOrder
}
