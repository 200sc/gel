package gel

import "math"

// A Vertex is a point in 3D space with an x, y, and z coordinate.
type Vertex struct {
	x, y, z float64
}

// Sub subtracts v2 from v and returns a Vertex of the difference.
func (v Vertex) Sub(v2 Vertex) Vertex {
	return Vertex{v.x - v2.x, v.y - v2.y, v.z - v2.z}
}

// Cross returns a Vertex of the cross product between v and v2
func (v Vertex) Cross(v2 Vertex) Vertex {
	return Vertex{v.y*v2.z - v.z*v2.y, v.z*v2.x - v.x*v2.z, v.x*v2.y - v.y*v2.x}
}

// Mul returns this Vertex with each of it's x,y, and z multiplied by n
func (v Vertex) Mul(n float64) Vertex {
	return Vertex{v.x * n, v.y * n, v.z * n}
}

// Dot calculates the dot product of two vertices
func (v Vertex) Dot(v2 Vertex) float64 {
	return v.x*v2.x + v.y*v2.y + v.z*v2.z
}

// Len returns the length or magnitude of this vertex as a vector
func (v Vertex) Len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// Unit converts this Vertex into a unit vector, by dividing it by it's magnitude
func (v Vertex) Unit() Vertex {
	return v.Mul(1.0 / v.Len())
}

// VMaxLen returns the maximum length from a set of vertices
func VMaxLen(vsv []Vertex) (max float64) {
	for _, v := range vsv {
		if v.Len() > max {
			max = v.Len()
		}
	}
	return
}
