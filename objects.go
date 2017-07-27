package gologo

import (
    "github.com/go-gl/mathgl/mgl32"
    "github.com/pkg/errors"
)

type Object struct {
    Config *ObjectTypeConfig
    ObjectType *ObjectType
    Model mgl32.Mat4
    ZOrder int
    Creation float64
    RenderData interface{}
}

type ObjectSet map[*Object]bool

func (o *Object) GetPrimitive() Primitive {
    return o.ObjectType.Primitive
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

func InitObjects(objects []*Object, objectTypes map[string]*ObjectType) error {
    for _, object := range objects {
        config := object.Config
        if config == nil {
            return errors.Errorf("Missing config for object: %v\n", object)
        }

        objectType := objectTypes[object.Config.Name]
        if objectType == nil {
            return errors.Errorf("Unknown object type specified: %v\n",
                object.Config.Name)
        }

        object.ObjectType = objectType
    }

    return nil
}
