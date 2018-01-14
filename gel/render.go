package gel

import (
	"image"
	"image/draw"
	"math"
	"os"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

type Render struct {
	*render.Sprite
	zbuff     []float64
	fdif      *image.RGBA
	tv        []Triangle
	tt        []Triangle
	tn        []Triangle
	w, h      int
	lastmouse mouse.Event
}

func NewRender(objfile, bmpfile string) (*Render, error) {
	fobj, err := os.Open(objfile)
	if err != nil {
		return nil, err
	}
	sp, err := render.LoadSprite(bmpfile)
	if err != nil {
		return nil, err
	}
	obj := oparse(fobj)
	return &Render{
		Sprite: sp,
		tv:     obj.Tvgen(),
		tt:     obj.Ttgen(),
		tn:     obj.Tngen(),
	}, nil
}

func (r *Render) Draw(buff draw.Image) {
	r.DrawOffset(buff, 0, 0)
}

func (r *Render) DrawOffset(buff draw.Image, xOff, yOff float64) {
	if mouse.LastEvent != r.lastmouse {
		mouseXt := mouse.LastEvent.X() * .005
		mouseYt := mouse.LastEvent.Y() * .005
		zbuff := make([]float64, r.w*r.h)
		ctr := Vertex{0.0, 0.0, 0.0}
		ups := Vertex{0.0, 1.0, 0.0}
		eye := Vertex{math.Sin(mouseXt), math.Sin(mouseYt), math.Sin(mouseXt)}
		z := eye.Sub(ctr).Unit()
		x := ups.Cross(z).Unit()
		y := z.Cross(x)
		for i := 0; i < len(r.tv); i++ {
			nrm := r.tn[i].ViewNrm(x, y, z)
			tex := r.tt[i]
			tri := r.tv[i].ViewTri(x, y, z, eye)
			per := tri.Perspective()
			vew := per.Viewport(floatgeom.Point2{float64(r.w), float64(r.h)})
			targ := Target{vew, nrm, tex, r.fdif}
			TDraw(r.h, zbuff, targ)
		}
		r.Sprite.DrawOffset(buff, xOff, yOff)
	}
	r.lastmouse = mouse.LastEvent
}
