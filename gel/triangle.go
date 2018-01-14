package gel

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/alg/floatgeom"
)

type Triangle struct {
	a, b, c Vertex
}

func (t Triangle) Viewport(field floatgeom.Point2) Triangle {
	w := field.Y() / 1.5
	h := field.Y() / 1.5
	x := field.X() / 2.0
	y := field.Y() / 4.0
	return Triangle{
		Vertex{w*t.a.x + x, h*t.a.y + y, (t.a.z + 1.0) / 1.5},
		Vertex{w*t.b.x + x, h*t.b.y + y, (t.b.z + 1.0) / 1.5},
		Vertex{w*t.c.x + x, h*t.c.y + y, (t.c.z + 1.0) / 1.5},
	}
}

func (t Triangle) Perspective() Triangle {
	c := 3.0
	za := 1.0 - t.a.z/c
	zb := 1.0 - t.b.z/c
	zc := 1.0 - t.c.z/c
	return Triangle{
		Vertex{t.a.x / za, t.a.y / za, t.a.z / za},
		Vertex{t.b.x / zb, t.b.y / zb, t.b.z / zb},
		Vertex{t.c.x / zc, t.c.y / zc, t.c.z / zc},
	}
}

func (t Triangle) BaryCenter(x, y int) Vertex {
	p := Vertex{float64(x), float64(y), 0.0}
	v0 := t.b.Sub(t.a)
	v1 := t.c.Sub(t.a)
	v2 := p.Sub(t.a)
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	v := (d11*d20 - d01*d21) / (d00*d11 - d01*d01)
	w := (d00*d21 - d01*d20) / (d00*d11 - d01*d01)
	u := 1.0 - v - w
	return Vertex{v, w, u}
}

func (t Triangle) Unit() Triangle {
	return Triangle{t.a.Unit(), t.b.Unit(), t.c.Unit()}
}

func (t Triangle) Mul(f float64) Triangle {
	return Triangle{t.a.Mul(f), t.b.Mul(f), t.c.Mul(f)}
}

func (t Triangle) ViewTri(x, y, z, eye Vertex) Triangle {
	return Triangle{
		Vertex{t.a.Dot(x) - x.Dot(eye), t.a.Dot(y) - y.Dot(eye), t.a.Dot(z) - z.Dot(eye)},
		Vertex{t.b.Dot(x) - x.Dot(eye), t.b.Dot(y) - y.Dot(eye), t.b.Dot(z) - z.Dot(eye)},
		Vertex{t.c.Dot(x) - x.Dot(eye), t.c.Dot(y) - y.Dot(eye), t.c.Dot(z) - z.Dot(eye)},
	}
}

func (t Triangle) ViewNrm(x, y, z Vertex) Triangle {
	return Triangle{
		Vertex{t.a.Dot(x), t.a.Dot(y), t.a.Dot(z)},
		Vertex{t.b.Dot(x), t.b.Dot(y), t.b.Dot(z)},
		Vertex{t.c.Dot(x), t.c.Dot(y), t.c.Dot(z)},
	}.Unit()
}

func TDraw(buff *image.RGBA, zbuff [][]float64, t Target) {
	x0 := int(math.Min(t.vew.a.x, math.Min(t.vew.b.x, t.vew.c.x)))
	y0 := int(math.Min(t.vew.a.y, math.Min(t.vew.b.y, t.vew.c.y)))
	x1 := int(math.Max(t.vew.a.x, math.Max(t.vew.b.x, t.vew.c.x)))
	y1 := int(math.Max(t.vew.a.y, math.Max(t.vew.b.y, t.vew.c.y)))
	dims := t.fdif.Bounds()
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			// Coordinate system is upwards.
			bc := t.vew.BaryCenter(x, y)
			if bc.x >= 0.0 && bc.y >= 0.0 && bc.z >= 0.0 {
				// But everything else here is rotated 90 degrees to accomodate a fast render cache.
				z := bc.x*t.vew.b.z + bc.y*t.vew.c.z + bc.z*t.vew.a.z
				if z > zbuff[x][y] {
					light := Vertex{0.0, 0.0, 1.0}
					varying := Vertex{light.Dot(t.nrm.b), light.Dot(t.nrm.c), light.Dot(t.nrm.a)}

					xx := (float64(dims.Max.X) - 1) * (0.0 + (bc.x*t.tex.b.x + bc.y*t.tex.c.x + bc.z*t.tex.a.x))
					yy := (float64(dims.Max.Y) - 1) * (1.0 - (bc.x*t.tex.b.y + bc.y*t.tex.c.y + bc.z*t.tex.a.y))
					intensity := bc.Dot(varying)
					var shading uint32
					if intensity > 0.0 {
						shading = uint32(intensity * 0xFF)
					}
					// Again, notice the rotated renderer (destination) but right side up image (source).
					zbuff[x][y] = z
					buff.Set(x, y, PShade(t.fdif.At(int(xx), int(yy)), shading))
				}
			}
		}
	}
}

func PShade(pixel color.Color, shading uint32) color.RGBA {
	r, g, b, a := pixel.RGBA()
	r /= 257
	r *= shading
	r >>= 0x08
	g /= 257
	g *= shading
	g >>= 0x08
	b /= 257
	b *= shading
	b >>= 0x08
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
