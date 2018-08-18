package gologo

import (
	"testing"

	"reflect"

	"github.com/go-gl/mathgl/mgl32"
)

var createTOTests = []struct {
	templateType string
	posX, posY   float32
	primitiveExp Primitive
	rendererExp  Renderer
}{
	{"SIMPLE_OBJECT", 200, 200, nil, nil},
	{"PRIM_OBJECT", 200, 200, &Circle{Radius: 50}, nil},
}

// TestCreateTemplateObject : Test basic Object creation by building the object
// from a template then reading values back out to confirm
// A simple object has no renderer
func TestCreateTemplateObject(t *testing.T) {
	var position mgl32.Vec3
	var obj *Object
	var prim Primitive
	var rend Renderer
	var err error
	var age int
	var x, y float32

	LoadObjectTemplates("testdata" + pathSeparator + "res")

	for _, tt := range createTOTests {
		position = mgl32.Vec3{tt.posX, tt.posY, 0.0}
		obj, err = CreateTemplateObject(tt.templateType, position)
		if err != nil {
			t.Errorf("Obj (%v) create failed with error: %v", tt.templateType, err)
		}

		// Test object has age
		age = obj.GetAge()
		if age < 0 {
			t.Errorf("Age was (%v) should not be negative", age)
		}

		// Test position
		x, y = obj.GetPosition()
		if x != tt.posX || y != tt.posY {
			t.Errorf("Location was x (%v), y (%v) should be x (%v), y (%v)",
				x, y, tt.posX, tt.posY)
		}

		// Check primitive
		prim = obj.GetPrimitive()

		if prim != tt.primitiveExp {
			Info.Printf("prim type is (%v)", reflect.TypeOf(prim))
			Info.Printf("exp type is (%v)", reflect.TypeOf(tt.primitiveExp))
			t.Errorf("Primitive is (%+v) should be (%+v)", prim, tt.primitiveExp)
		}

		// Check renderer
		rend = obj.GetRenderer()
		if rend != tt.rendererExp {
			t.Errorf("Renderer is (%+v) should be (%+v)", rend, tt.rendererExp)
		}
	}
}
