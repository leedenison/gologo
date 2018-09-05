package gologo

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/opengl"
	"github.com/leedenison/gologo/render"
	"github.com/leedenison/gologo/time"
)

/////////////////////////////////////////////////////////////
// TextRenderer
//

type TextRenderer struct {
	MeshRenderers map[byte]render.Renderer
	CharWidths    map[byte]float32
	CharSpacer    float32
	Text          []byte
	Transforms    []mgl32.Mat4
}

func (r *TextRenderer) Render(model mgl32.Mat4) {
	r.RenderAt(model, map[int]interface{}{})
}

func (r *TextRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	if len(r.Transforms) != len(r.Text) {
		r.InitTransforms()
	}

	for i := 0; i < len(r.Text); i++ {
		renderer, ok := r.MeshRenderers[r.Text[i]]

		if ok {
			renderer.RenderAt(
				model.Mul4(r.Transforms[i]),
				map[int]interface{}{})
		}
	}
}

func (r *TextRenderer) InitTransforms() {
	var translate float32

	if len(r.Text) == 0 {
		return
	}

	r.Transforms = make([]mgl32.Mat4, len(r.Text))
	for count := 0; count < len(r.Text); count++ {
		width, ok := r.CharWidths[r.Text[count]]
		if ok {
			r.Transforms[count] = mgl32.Translate3D(
				translate+(width/2), 0, 0)
			translate += width + r.CharSpacer
		}
	}

	for count := 0; count < len(r.Text); count++ {
		r.Transforms[count] = mgl32.Translate3D(-translate/2, 0, 0).
			Mul4(r.Transforms[count])
	}
}

func (r *TextRenderer) DebugRender(model mgl32.Mat4) {
	r.DebugRenderAt(model, map[int]interface{}{})
}

func (r *TextRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	if len(r.Transforms) != len(r.Text) {
		r.InitTransforms()
	}

	for i := 0; i < len(r.Text); i++ {
		renderer, ok := r.MeshRenderers[r.Text[i]]

		if ok {
			renderer.DebugRenderAt(
				model.Mul4(r.Transforms[i]),
				map[int]interface{}{})
		}
	}
}

func (r *TextRenderer) Animate(model mgl32.Mat4) {}

func (r *TextRenderer) Clone() render.Renderer {
	return &TextRenderer{
		MeshRenderers: r.MeshRenderers,
		CharSpacer:    r.CharSpacer,
	}
}

/////////////////////////////////////////////////////////////
// ExplosionRenderer
//

type Particle struct {
	Model    mgl32.Mat4
	Velocity mgl32.Mat4
	Age      float64
	Renderer render.Renderer
}

type ExplosionRenderer struct {
	Renderers     []render.Renderer
	ParticleCount int
	Particles     []*Particle
	MaxAge        float64
}

func (r *ExplosionRenderer) Render(model mgl32.Mat4) {
	r.RenderAt(model, map[int]interface{}{})
}

func (r *ExplosionRenderer) RenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	for i := 0; i < len(r.Particles); i++ {
		age := float64(time.GetTickTime()) - r.Particles[i].Age
		r.Particles[i].Renderer.RenderAt(
			model.Mul4(r.Particles[i].Model),
			map[int]interface{}{
				opengl.UniformAlpha: 1.0 - age/r.MaxAge,
			})
	}
}

func (r *ExplosionRenderer) DebugRender(model mgl32.Mat4) {
	r.DebugRenderAt(model, map[int]interface{}{})
}

func (r *ExplosionRenderer) DebugRenderAt(model mgl32.Mat4, custom map[int]interface{}) {
	log.Trace.Printf("ExplosionRenderer: Model matrix:\n%v\n", model)
}

func (r *ExplosionRenderer) Animate(model mgl32.Mat4) {
	if r.Particles == nil {
		r.Particles = make([]*Particle, r.ParticleCount)
		for i := 0; i < r.ParticleCount; i++ {
			r.Particles[i] = r.createRandomParticle()
		}
	}

	for j := 0; j < len(r.Particles); j++ {
		age := float64(time.GetTickTime()) - r.Particles[j].Age
		if age >= r.MaxAge {
			if len(r.Particles) > 1 {
				r.Particles = append(r.Particles[:j], r.Particles[j+1:]...)
			} else {
				r.Particles = r.Particles[0:0]
			}
			j--
		} else {
			r.Particles[j].Model = r.Particles[j].Velocity.Mul4(r.Particles[j].Model)
		}
	}
}

func (r *ExplosionRenderer) createRandomParticle() *Particle {
	renderer := rand.Intn(len(r.Renderers))
	velocity := mgl32.HomogRotate3DZ(rand.Float32() * 0.01)
	velocity = velocity.Mul4(mgl32.Translate3D(rand.Float32(), rand.Float32(), 0.0))

	return &Particle{
		Renderer: r.Renderers[renderer],
		Velocity: velocity,
		Age:      float64(time.GetTickTime()),
		Model:    mgl32.Ident4(),
	}
}

func (r *ExplosionRenderer) Clone() render.Renderer {
	return &ExplosionRenderer{
		Renderers: r.Renderers,
		MaxAge:    r.MaxAge,
	}
}
