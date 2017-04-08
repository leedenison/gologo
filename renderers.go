package gologo

import (
    "github.com/AllenDang/w32"
)

var CIRCLE_RENDERER = CircleRenderer {
    Painter: painters[PAINTER_OBJ],
}

type Renderer interface {
    Render(Object, w32.HDC, *RenderState)
}

type CircleRenderer struct {
    Painter *Painter
}

func (r *CircleRenderer) Render(obj Object, hdc w32.HDC, state *RenderState) {
    switch t := obj.(type) {
    case *Circle:
        if r.Painter.Type != OBJECT_EMPTY {
            SelectPainterIfNeeded(hdc, r.Painter, state)

            w32.Ellipse(
                hdc,
                int(t.Particle.Position.x - t.Radius),
                int(windowHeight - int32(t.Particle.Position.y - t.Radius)),
                int(t.Particle.Position.x + t.Radius),
                int(windowHeight - int32(t.Particle.Position.y + t.Radius)))
        }
    default:
        return
    }
}