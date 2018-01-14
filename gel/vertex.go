package gel

import "math"

type Vertex struct {
	x, y, z float64
}

func (v Vertex) Sub(v2 Vertex) Vertex {
	return Vertex{v.x - v2.x, v.y - v2.y, v.z - v2.z}
}

func (v Vertex) Cross(v2 Vertex) Vertex {
	return Vertex{v.y*v2.z - v.z*v2.y, v.z*v2.x - v.x*v2.z, v.x*v2.y - v.y*v2.x}
}

func (v Vertex) Mul(n float64) Vertex {
	return Vertex{v.x * n, v.y * n, v.z * n}
}

func (v Vertex) Dot(v2 Vertex) float64 {
	return v.x*v2.x + v.y*v2.y + v.z*v2.z
}

func (v Vertex) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vertex) Unit() Vertex {
	return v.Mul(1.0 / v.Len())
}

func VMaxLen(vsv []Vertex) (max float64) {
	for _, v := range vsv {
		if v.Len() > max {
			max = v.Len()
		}
	}
	return
}
