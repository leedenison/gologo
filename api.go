package gologo

func CreateCircle(x, y int32) {
    objects[next_object_id] = &Circle{
        Id: next_object_id,
        Centre: Vector{ x: float64(x), y: float64(y) },
        Radius: 20,
        Renderer: renderers[RENDER_OBJ],
    }
    next_object_id++
}

func CreateCircleWithSpeed(x, y, speedx, speedy int32) {
    objects[next_object_id] = &Circle{
        Id: next_object_id,
        Centre: Vector{ x: float64(x), y: float64(y) },
        Radius: 20,
        Renderer: renderers[RENDER_OBJ],
    }
    movables[next_object_id] = Vector{ x: float64(speedx), y: float64(speedy) }
    next_object_id++
}