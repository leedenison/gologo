package gologo

func CreateCircle(x, y int32) {
    objects[nextObjectId] = &Circle{
        Id: nextObjectId,
        Centre: Vector{ x: float64(x), y: float64(y) },
        Radius: 20,
        Renderer: renderers[RENDER_OBJ],
    }
    nextObjectId++
}

func CreateCircleWithSpeed(x, y, speedx, speedy int32) {
    objects[nextObjectId] = &Circle{
        Id: nextObjectId,
        Centre: Vector{ x: float64(x), y: float64(y) },
        Radius: 20,
        Renderer: renderers[RENDER_OBJ],
    }
    movables[nextObjectId] = Vector{
        x: float64(speedx) / SPEED_MULT,
        y: float64(-speedy) / SPEED_MULT,
    }
    nextObjectId++
}

func SetGravity(x, y int32) {
    gravity.x = float64(x) / SPEED_MULT
    gravity.y = float64(y) / SPEED_MULT
}
