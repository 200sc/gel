package gel

import (
	"bufio"
)

type Vertex struct {
    x,y,z float64
}

type Vertices []Vertex

type Face struct {
    va, vb, vc int
    ta, tb, tc int
    na, nb, nc int
}

type Faces []Face

type Obj struct {
    vsv, vsn, vst Vertices
    fs Faces 
}

type Triangle struct {
    a, b, c Vertex
}

type Triangles []Triangle

type Target struct {
    vew, nrm, tex Triangle
    fdif *SDL_Surface
}

type Sdl struct {
    window *SDL_Window
    renderer *SDL_Renderer
    canvas *SDL_Texture
    xres, yres int 
}

type Input struct {
    xt, yt, sens float64
    key *uint8
}

func (v Vertex) Sub(v2 Vertex) Vertex {
    return Vertex{v.x - v2.x, v.y-v2.y, v.z-v2.z}
}

func (v Vertex) Cross(v2 Vertex) Vertex {
    return Vertex{v.y * v2.z - v.z * v2.y, v.z * v2.x - v.x * v2.z, v.x * v2.y - v.y * v2.x }
}

func (v Vertex) Mul(n float64) {
    return Vertex{v.x * n, v.y * n, v.z * n}
}

func (v Vertex) Dot(v2 Vertex) float64 {
    return v.x * v2.x + v.y * v2.y + v.z * v2.z
}

func (v Vertex) Len() float64 {
    return math.Sqrt(v.x * v.x + v.y * v.y + v.z * v.z)
}

func (v Vertex) Unit() Vertex {
    return v.Mul(1.0 / v.Len())
}

func (t Triangle) Unit() Triangle {
    return Triangle{t.a.Unit(), t.b.Unit(), t.c.Unit()}
}

func (t Triangle) Mul(f float64) Triangle {
    return Triangle{t.a.Mul(f), t.b.Mul(f), t.c.Mul(f)}
}

func VMaxLen(vsv []Vertex) (max float64) {
    for _, v := range vsv {
        if v.Len() > max {
            max = v.Len()
        }
    }
    return
}

func (o Obj) Tvgen() []Triangle {
    scale := 1.0 / VMaxLen(obj.vsv)
    tv := make([]Triangle, obj.fs.count)
    for i := 0; i < obj.fs.count; i++ {
        tv[i] = Triangle{
            obj.vsv.vertex[obj.fs.face[i].va],
            obj.vsv.vertex[obj.fs.face[i].vb],
            obj.vsv.vertex[obj.fs.face[i].vc],
        }.Mul(scale)
    }
    return tv
}

func (o Obj) Tngen() []Triangle {
    scale := 1.0 / VMaxLen(obj.vsv)
    tn := make([]Triangle, obj.fs.count)
    for i := 0; i < obj.fs.count; i++ {
        tv[i] = Triangle{
            obj.vsv.vertex[obj.fs.face[i].na],
            obj.vsv.vertex[obj.fs.face[i].nb],
            obj.vsv.vertex[obj.fs.face[i].nc],
        }.Mul(scale)
    }
    return tn
}

func (o Obj) Ttgen() []Triangle {
    scale := 1.0 / VMaxLen(obj.vsv)
    tt := make([]Triangle, obj.fs.count)
    for i := 0; i < obj.fs.count; i++ {
        tv[i] = Triangle{
            obj.vsv.vertex[obj.fs.face[i].ta],
            obj.vsv.vertex[obj.fs.face[i].tb],
            obj.vsv.vertex[obj.fs.face[i].tc],
        }.Mul(scale)
    }
    return tt
}

func (t Triangle) Viewport(sdl Sdl) Triangle {
    w := sdl.yres / 1.5
    h := sdl.yres / 1.5
    x := sdl.xres / 2.0
    y := sdl.yres / 4.0
    return Triangle{
        { w * t.a.x + x, h * t.a.y + y, (t.a.z + 1.0) / 1.5 },
        { w * t.b.x + x, h * t.b.y + y, (t.b.z + 1.0) / 1.5 },
        { w * t.c.x + x, h * t.c.y + y, (t.c.z + 1.0) / 1.5 },
    }
}

func (t Triangle) Perpspective() Triangle {
    c := 3.0;
    za := 1.0 - t.a.z / c
    zb := 1.0 - t.b.z / c
    zc := 1.0 - t.c.z / c
    return Triangle{
        { t.a.x / za, t.a.y / za, t.a.z / za },
        { t.b.x / zb, t.b.y / zb, t.b.z / zb },
        { t.c.x / zc, t.c.y / zc, t.c.z / zc },
    }
}

func (t Triangle) BaryCenter(x, y int) Vertex {
      p := Vertex{ float64(x), float64(y), 0.0 }
      v0 := t.b.Sub(t.a)
      v1 := t.c.Sub(t.a)
      v2 := p.Sub(t.a)
      d00 := v0.Dot(v0)
      d01 := v0.Dot(v1)
      d11 := v1.Dot(v1)
      d20 := v2.Dot(v0)
      d21 := v2.Dot(v1)
      v := (d11 * d20 - d01 * d21) / (d00 * d11 - d01 * d01)
      w := (d00 * d21 - d01 * d20) / (d00 * d11 - d01 * d01)
      u = 1.0 - v - w;
    return Vertex{v, w, u}
}

// todo: use a color.RGBA type or something instead of a uint32
func PShade(pixel uint32, shading int) uint32 {
    r := (((pixel >> 0x10) /****/) * shading) >> 0x08;
    g := (((pixel >> 0x08) & 0xFF) * shading) >> 0x08;
    b := (((pixel /*****/) & 0xFF) * shading) >> 0x08;
    return r << 0x10 | g << 0x08 | b;
}

func TDraw(yres int, pixel uint32, zbuff []float64, t Target) {
     x0 := int(math.Min(t.vew.a.x, math.Min(t.vew.b.x, t.vew.c.x)))
     y0 := int(math.Min(t.vew.a.y, math.Min(t.vew.b.y, t.vew.c.y)))
     x1 := int(math.Max(t.vew.a.x, math.Max(t.vew.b.x, t.vew.c.x)))
     y1 := int(math.Max(t.vew.a.y, math.Max(t.vew.b.y, t.vew.c.y)))
     for x := x0; x <= x1; x++ {
         for y := y0; y <= y1; y++ {
            // Coordinate system is upwards. 
            bc := t.vew.BaryCenter(x,y)
            if bc.x >= 0.0 && bc.y >= 0.0 && bc.z >= 0.0 {
                // But everything else here is rotated 90 degrees to accomodate a fast render cache.
                z := bc.x * t.vew.b.z + bc.y * t.vew.c.z + bc.z * t.vew.a.z
                if z > zbuff[y + x * yres] {
                    light := Vertex{ 0.0, 0.0, 1.0 }
                    varying := Vertex{ light.Dot(t.nrm.b), light.Dot(t.nrm.c), light.Dot(t.nrm.a) }
                    pixels := t.fdif.pixels
                    xx := (t.fdif.w - 1) * (0.0 + (bc.x * t.tex.b.x + bc.y * t.tex.c.x + bc.z * t.tex.a.x))
                    yy := (t.fdif.h - 1) * (1.0 - (bc.x * t.tex.b.y + bc.y * t.tex.c.y + bc.z * t.tex.a.y))
                    intensity := bc.Dot(varying)
                    shading := 0.0
                    if intensity > 0.0 {
                        shading = intensity * 0xFF
                    }
                    // Again, notice the rotated renderer (destination) but right side up image (source).
                    zbuff[y + x * yres] = z
                    pixel[y + x * yres] = PShade(pixels[xx + yy * t.fdif.w], shading)
                }
            }
        }
    }
}

func (t Triangle) ViewTri(x,y,z,eye Vertex) Triangle {
    return Triangle{
        { t.a.Dot(x) - x.Dot(eye), t.a.Dot(y) - y.Dot(eye), t.a.Dot(z) - z.Dot(eye) },
        { t.b.Dot(x) - x.Dot(eye), t.b.Dot(y) - y.Dot(eye), t.b.Dot(z) - z.Dot(eye) },
        { t.c.Dot(x) - x.Dot(eye), t.c.Dot(y) - y.Dot(eye), t.c.Dot(z) - z.Dot(eye) },
    }
}

func (t Triangle) ViewNrm(n Triangle, x,y,z Vertex) Triangle {
    return Triangle{
        { n.a.Dot(x), n.a.Dot(y), n.a.Dot(z) },
        { n.b.Dot(x), n.b.Dot(y), n.b.Dot(z) },
        { n.c.Dot(x), n.c.Dot(y), n.c.Dot(z) },
    }.Unit()
}

func oparse(f os.File) {
    size := 128
    vsv := make([]Vertex, 0, size)
    vsn := make([]Vertex, 0, size)
    vst := make([]Vertex, 0, size)
    fs := make([]Face, 0, size)
    
    scn := bufio.NewScanner(f)
    defer func() { 
        err := f.Close()
        if err != nil{
            dlog.Error(err)
        }
    }()

    for scn.Scan() {
        var f Face
        var v Vertex

        line := scn.Text()

        if line[0] == 'v' && line[1] == 'n' {
            fmt.Sscanf(line, "vn %f %f %f", &v.x, &v.y, &v.z)
            vsn.count++
            vsn.vertex[vsn.count] = v
        } else if line[0] == 'v' && line[1] == 't' {
            fmt.Sscanf(line, "vt %f %f %f", &v.x, &v.y, &v.z)
            vsv.count++
            vst.vertex[vst.count] = v
        } else if (line[0] == 'v') {
            fmt.Sscanf(line, "v %f %f %f", &v.x, &v.y, &v.z)
            vsv.count++
            vsv.vertex[vsv.count] = v
        } else if (line[0] == 'f') {
            fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &f.va, &f.ta, &f.na, &f.vb, &f.tb, &f.nb, &f.vc, &f.tc, &f.nc)
            fs.count++
            fs.face[fs.count] = Face{
                f.va - 1, f.vb - 1, f.vc - 1,
                f.ta - 1, f.tb - 1, f.tc - 1,
                f.na - 1, f.nb - 1, f.nc - 1,
            }
        }
    }
    return Obj{vsv, vsn, vst, fs}
}

func main() {
    r, err := NewRender("model/salesman.obj", "model/salesman.bmp")
    if err != nil {
        dlog.Error(err)
        return
    }
    render.Draw(r)
    oak.Init(800, 600)
}

type Render struct {
    *render.Sprite
    zbuff []float64
    fdif //???
    obj Object
    tv []Triangle
    tt []Triangle
    tn []Triangle
}

func NewRender(objfile, bmpfile string) (*Render, error) {
    fobj, err := os.Open("model/salesman.obj")
    if err != nil {
        return nil, err
    }
    // ???
    fdif = sload("model/salesman.bmp")
    obj = oparse(fobj)
    return &Render{
        tv: tvgen(obj),
        tt: ttgen(obj),
        tn: tngen(obj),
        obj: obj,
        fdid: fdif,
    }, nil
}

func (r *Render) Draw(buff draw.Image) {
    r.DrawOffset(buff, 0, 0)
}

func (r *Render) DrawOffset(buff draw.Image, xOff, yOff float64) {
    // todo: move most of main here
    zbuff = make([]float64, sdl.xres * sdl.yres)
    for {
        uint32_t* const pixel = slock(sdl)
        zbuff = make([]float64, sdl.xres * sdl.yres)
        pixel = make([]uint32, sdl.xres * sdl.yres)
        ctr := Vertex{ 0.0, 0.0, 0.0 }
        ups := Vertex{ 0.0, 1.0, 0.0 }
        eye := Vertex{ sinf(input.xt), sinf(input.yt), cosf(input.xt) }
        z = eye.Sub(ctr).Unit()
        x = ups.Cross(z).Unit()
        y = z.Cross(x)
        for i := 0; i < len(tv); i++ {
            nrm := tviewnrm(tn.triangle[i], x, y, z)
            tex := tt.triangle[i]
            tri := tviewtri(tv.triangle[i], x, y, z, eye)
            per := tpersp(tri)
            vew := tviewport(per, sdl)
            targ := Target{ vew, nrm, tex, fdif }
            tdraw(sdl.yres, pixel, zbuff, targ)
        }
    }
    r.Sprite.DrawOffset(buff, xOff, yOff)
}