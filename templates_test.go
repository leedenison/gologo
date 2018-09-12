package gologo

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-test/deep"
	"github.com/leedenison/gologo/mocks"
	"github.com/leedenison/gologo/opengl"
	"github.com/leedenison/gologo/render"
)

var createTOTests = []struct {
	name, templateType string
	posX, posY         float32
	primitiveExp       Primitive
	rendererExp        render.Renderer
	errExp             error
}{
	{"simple object", "SIMPLE_OBJECT", 200, 300, nil, nil, nil},
	{"defined circle", "CIRCLE_OBJECT", 200, 300, &Circle{Radius: 50}, nil, nil},
	{"undefined circle and renderer", "UNDEF_CIRCLE_REND_OBJECT", 200, 300,
		&Circle{Radius: 65},
		&opengl.MeshRenderer{Shader: nil, Mesh: 0, Uniforms: nil,
			MeshVertices: []float32{-100, -100, 0, 0, 1, 100, 100, 0, 1, 0,
				-100, 100, 0, 0, 0, -100, -100, 0, 0, 1,
				100, -100, 0, 1, 1, 100, 100, 0, 1, 0},
			VertexCount: 6}, nil},
}

// TestCreateTemplateObject : Test Object creation by building the object
// from a template then reading values back out to confirm
func TestCreateTemplateObject(t *testing.T) {
	var position mgl32.Vec3
	var obj *Object
	var prim Primitive
	var rend render.Renderer
	var err error
	var x, y float32
	var diff []string

	LoadObjectTemplates("testdata" + pathSeparator + "res")

	opengl.CreateMeshRenderer = mocks.CreateVerticesOnlyMeshRendererImpl

	for _, tc := range createTOTests {
		t.Run(tc.name, func(t *testing.T) {
			position = mgl32.Vec3{tc.posX, tc.posY, 0.0}
			obj, err = CreateTemplateObject(tc.templateType, position)
			if err != tc.errExp {
				t.Errorf("Create object error was (%v) error should have been (%v)", err, tc.errExp)
			}

			// Test position
			x, y = obj.GetPosition()
			if x != tc.posX || y != tc.posY {
				t.Errorf("Location was x (%v), y (%v) should be x (%v), y (%v)",
					x, y, tc.posX, tc.posY)
			}

			// Check primitive
			prim = obj.GetPrimitive()
			if !reflect.DeepEqual(prim, tc.primitiveExp) {
				t.Errorf("Primitive is (%+v) should be (%+v)", prim, tc.primitiveExp)
			}

			// Check renderer
			rend = obj.GetRenderer()

			diff = deep.Equal(rend, tc.rendererExp)
			if diff != nil {
				t.Error("Generated renderer is different:", diff)
			}
		})
	}
}
