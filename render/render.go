package render

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer interface {
	Render(model mgl32.Mat4)
	RenderAt(model mgl32.Mat4, uniforms map[int]interface{})
	DebugRender(model mgl32.Mat4)
	DebugRenderAt(model mgl32.Mat4, uniforms map[int]interface{})
	Animate(model mgl32.Mat4)
	Clone() Renderer
}
