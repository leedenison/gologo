package gologo

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

// TestObjectCreation : Test basic Object creation by building the object
// then reading values back out to confirm values
func TestObjectCreation(t *testing.T) {
	// From here only use public functions
	defer Cleanup()
	Init()

	LoadObjectTemplates()

	objBuilder := Builder()
	objBuilder.
		SetPosition(200, 200).
		Build("BIG_BLUE_SQUARE")

	model := mgl32.Translate3D(400, 400, 0.0).Mul4(mgl32.Ident4().Mul4(mgl32.Ident4()))

	obj := CreateObject(model)
	obj.SetZOrder(1)

	// Create meshRenderer
	meshRenderer, err := CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{uniformColor},
		map[int]interface{}{uniformColor: mgl32.Vec4{1, 0, 0, 1}},
		[]float32{
			-50, -50, 0, 0, 1,
			50, 50, 0, 1, 0,
			-50, 50, 0, 0, 0,
			-50, -50, 0, 0, 1,
			50, -50, 0, 1, 1,
			50, 50, 0, 1, 0,
		},
	)

	if err != nil {
		errors.Wrapf(err, "Failed to create defaulted MeshRenderer (%v)", meshRenderer)
		t.Error(err)
	}

	obj.SetRenderer(meshRenderer, true)

	TagRender(obj)

	Run()
}
