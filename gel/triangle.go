package gel

import (
	"image"
	"image/color"
	"math"

	"github.com/oakmound/oak/alg/floatgeom"
)

// A Triangle is a set of three points or vertices.
type Triangle struct {
	a, b, c Vertex
}

// Unit returns a Triangle where each vertex is converted into a unit vector.
func (t Triangle) Unit() Triangle {
	return Triangle{t.a.Unit(), t.b.Unit(), t.c.Unit()}
}

// Mul returns a Triangle with each vertex magnified by f.
func (t Triangle) Mul(f float64) Triangle {
	return Triangle{t.a.Mul(f), t.b.Mul(f), t.c.Mul(f)}
}

// More documentation needed in the rest of this file,
// where the significant math occurs

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

func TDraw(buff *image.RGBA, zbuff [][]float64, vew, nrm, tex Triangle, textureData *image.RGBA) {
	x0 := int(math.Min(vew.a.x, math.Min(vew.b.x, vew.c.x)))
	y0 := int(math.Min(vew.a.y, math.Min(vew.b.y, vew.c.y)))
	x1 := int(math.Max(vew.a.x, math.Max(vew.b.x, vew.c.x)))
	y1 := int(math.Max(vew.a.y, math.Max(vew.b.y, vew.c.y)))
	dims := textureData.Bounds()
	buffH := buff.Bounds().Max.Y
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			// Coordinate system is upwards.
			bc := vew.BaryCenter(x, y)
			if bc.x >= 0.0 && bc.y >= 0.0 && bc.z >= 0.0 {
				// But everything else here is rotated 90 degrees to accomodate a fast render cache.
				z := bc.x*vew.b.z + bc.y*vew.c.z + bc.z*vew.a.z
				if z > zbuff[x][y] {
					light := Vertex{0.0, 0.0, 1.0}
					varying := Vertex{light.Dot(nrm.b), light.Dot(nrm.c), light.Dot(nrm.a)}

					xx := (float64(dims.Max.X) - 1) * (0.0 + (bc.x*tex.b.x + bc.y*tex.c.x + bc.z*tex.a.x))
					yy := (float64(dims.Max.Y) - 1) * (1.0 - (bc.x*tex.b.y + bc.y*tex.c.y + bc.z*tex.a.y))
					intensity := bc.Dot(varying)
					var shading uint32
					if intensity > 0.0 {
						shading = uint32(intensity * 0xFF)
					}
					// Again, notice the rotated renderer (destination) but right side up image (source).
					zbuff[x][y] = z
					// Change from the original gel: we subtract y from buffH because,
					// somewhere, I (200sc) messed up the translation and we accidentally
					// are rendering everything upsidedown. This is the easiest fix!
					buff.Set(x, buffH-y, PShade(textureData.At(int(xx), int(yy)), shading))
				}
			}
		}
	}
}

// PShade takes a color and applies shading to it by magnifying it's rgb values
func PShade(pixel color.Color, shading uint32) color.RGBA {
	r, g, b, a := pixel.RGBA()
	// r,g, and b are divided by 257 because the .RGBA() function returns
	// values on a 16 bit scale instead of an 8 bit scale.
	// They are then magnified by shading, not overflowing because they are
	// int32 values, and shifted to the right to represent division after
	// shading.
	// Todo: shading should probably by a float64 instead of multiplying and
	// then shifting
	r = ((r / 257) * shading) >> 0x08
	g = ((g / 257) * shading) >> 0x08
	b = ((b / 257) * shading) >> 0x08
	// r,g, and b need to be cast back to uint8s to create a color.RGBA value
	// from them.
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
