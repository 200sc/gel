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
		Sprite: render.NewEmptySprite(0, 0, 800, 600),
		tv:     obj.Tvgen(),
		tt:     obj.Ttgen(),
		tn:     obj.Tngen(),
		fdif:   sp.GetRGBA(),
		w:      800,
		h:      600,
	}, nil
}

func (r *Render) Draw(buff draw.Image) {
	r.DrawOffset(buff, 0, 0)
}

func (r *Render) DrawOffset(buff draw.Image, xOff, yOff float64) {
	if mouse.LastEvent != r.lastmouse {
		r.Sprite.SetRGBA(image.NewRGBA(r.Sprite.GetRGBA().Bounds()))

		mouseXt := mouse.LastEvent.X() * .005
		mouseYt := mouse.LastEvent.Y() * .005
		zbuff := make([][]float64, r.w)
		for i := range zbuff {
			zbuff[i] = make([]float64, r.h)
			for j := range zbuff[i] {
				zbuff[i][j] = -math.MaxFloat64
			}
		}
		ctr := Vertex{0.0, 0.0, 0.0}
		ups := Vertex{0.0, 1.0, 0.0}
		eye := Vertex{math.Sin(mouseXt), math.Sin(mouseYt), math.Cos(mouseXt)}
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
			TDraw(r.Sprite.GetRGBA(), zbuff, targ)
		}
	}
	r.lastmouse = mouse.LastEvent
	r.Sprite.DrawOffset(buff, xOff, yOff)
}
