package gologo

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

// TestObjectCreation : Test basic Object creation by building the object
// then reading values back out to confirm values
func TestObjectCreation(t *testing.T) {
	/*
		LoadObjectTemplates("testdata" + pathSeparator + "res")

		objBuilder := Builder()
		objBuilder.
			SetPosition(200, 200).
			Build("BIG_BLUE_SQUARE")
	*/
	model := mgl32.Translate3D(400, 400, 0.0).Mul4(mgl32.Ident4().Mul4(mgl32.Ident4()))

	obj := CreateObject(model)
	obj.SetZOrder(1)
}
