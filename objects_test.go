package gologo

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

// TestObjectCreation : Test basic Object creation by building the object
// then reading values back out to confirm values
func TestObjectCreation(t *testing.T) {
	// From here only use public functions
	defer Cleanup()
	Init()

	model := mgl32.Translate3D(200, 200, 0.0).Mul4(mgl32.Ident4().Mul4(mgl32.Ident4()))

	obj := CreateObject(model)
	obj.SetZOrder(1)

	// Create meshRenderer
	meshRenderer := CreateMeshRenderer(
		"ORTHO_VERTEX_SHADER",
		"COLOR_FRAGMENT_SHADER",
		[]int{uniformColor},
		map[int]interface{ uniformColor }{mgl32.Vec4(1, 1, 1, 0)},
		[]float32{
			-1, -1, 0, 0, 1,
			1, 1, 0, 1, 0,
			-1, 1, 0, 0, 0,
			-1, -1, 0, 0, 1,
			1, -1, 0, 1, 1,
			1, 1, 0, 1, 0,
		},
	)

	Run()
	// From here directly query the object for correctness
	return
}
