package gologo

func CreateCircle(x, y int32) {
    objects[next_object_id] = &Circle{
        Centre: Vector{ x: float64(x), y: float64(y) },
        Radius: 20,
        Renderer: renderers[RENDER_OBJ],
    }
    movables[next_object_id] = Vector{ x: 0, y: 0 }
    next_object_id++
}

