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
