package gologo

type Matrix2 [4]float64

func (m *Matrix2) MultiplyVector2(v *Vector2) *Vector2 {
    return &Vector2 {
        x: m[0] * v.x + m[1] * v.y,
        y: m[2] * v.x + m[3] * v.y,
    }
}

func (m *Matrix2) ApplyTransform(v *Vector2) *Vector2 {
    return m.MultiplyVector2(v)
}

func (m *Matrix2) MultiplyMatrix2(n *Matrix2) *Matrix2 {
    return &Matrix2 {
        m[0] * n[0] + m[1] * n[2],
        m[0] * n[1] + m[1] * n[3],
        m[2] * n[0] + m[3] * n[2],
        m[2] * n[1] + m[3] * n[3],
    }
}

func (m *Matrix2) Determinant() float64 {
    return m[0] * m[3] - m[1] * m[2]
}

func (m *Matrix2) SetInverse(n *Matrix2) {
    determinant := n.Determinant()

    if (determinant == 0) {
        return
    }

    inverse := float64(1.0) / determinant

    // Store n[0] in case m == n
    n0 := n[0]

    m[0] = n[3] * inverse
    m[1] = -n[1] * inverse
    m[2] = -n[2] * inverse
    m[3] = n0 * inverse
}

func (m *Matrix2) Inverse() *Matrix2 {
    result := Matrix2 {}
    result.SetInverse(m)
    return &result
}

func (m *Matrix2) Invert() {
    m.SetInverse(m)
}

func (m *Matrix2) Transpose() {
    m1 := m[1]

    m[1] = m[2]
    m[2] = m1
}

func (m *Matrix2) SetOrientation(o *Vector2) {
    m[0] = o.x
    m[1] = -o.y
    m[2] = o.y
    m[3] = o.x
}