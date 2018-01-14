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

// The Render type implements oak/render.Renderable, letting it be drawn
// to an oak screen. Most of the functions of that interface are satisfied
// by embedding *render.Sprite into the struct, but we replace the Draw
// functions, which are called each frame, below.
type Render struct {
	// a render.Sprite has a position and a buffer of image data which
	// it uses to draw to the screen at that position.
	*render.Sprite
	// the textureData is the local texture file (.bmp in the original, .png in this version)
	// that is referred to to color each triangle face
	textureData *image.RGBA
	// tv, tt, and tn need better names
	tv []Triangle
	tt []Triangle
	tn []Triangle
	// the last mouse event is stored on the render so that when it changes,
	// the render knows to update what should be drawn.
	lastmouse mouse.Event
}

// NewRender creates a Render type to be drawn to screen.
// If it fails, it will return nil and the cause for failure.
// The inputs are the object and texture file paths to be loaded
// and the width and height of the render to be drawn.
func NewRender(objfile, texfile string, w, h int) (*Render, error) {
	// Open both the object and texture file
	// if either fails to load, return nil and
	// the cause.
	fobj, err := os.Open(objfile)
	if err != nil {
		return nil, err
	}
	tex, err := render.LoadSprite(texfile)
	if err != nil {
		return nil, err
	}
	// parse the object file into an Object struct
	obj := oparse(fobj)
	// return a Render built on that object struct and
	// the texture file.
	return &Render{
		Sprite:      render.NewEmptySprite(0, 0, w, h),
		tv:          obj.Tvgen(),
		tt:          obj.Ttgen(),
		tn:          obj.Tngen(),
		textureData: tex.GetRGBA(),
		// We set lastmouse to have the Reset event so that our
		// equality check in DrawOffset will fail on the first
		// draw frame. This allows for the render to be drawn
		// before any mouse input is recorded.
		lastmouse: mouse.Event{Event: "Reset"},
	}, nil
}

// Draw expects the render to draw itself to the input buffer.
func (r *Render) Draw(buff draw.Image) {
	// To avoid duplicating logic, Draw just calls DrawOffset with 0 offsets.
	r.DrawOffset(buff, 0, 0)
}

// DrawOffset expects the render to draw itself to the input buffer,
// offset from it's logical coordinates by xOff and yOff for x and y respectively
func (r *Render) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// If there hasn't been new mouse input, so the 3d model has not been rotated
	// since it was processed last, don't re-process the model's rotation.
	if mouse.LastEvent != r.lastmouse {
		// Reset the backing Sprite's color buffer to be empty, so we avoid
		// smearing with the last drawn frame.
		bounds := r.Sprite.GetRGBA().Bounds()
		r.Sprite.SetRGBA(image.NewRGBA(bounds))
		w := bounds.Max.X
		h := bounds.Max.Y

		// Get the mouse position and scale it down so we can use it for the
		// render's rotation.
		mouseXt := mouse.LastEvent.X() * .005
		mouseYt := mouse.LastEvent.Y() * .005
		// Set up a zbuffer so we know what pixels we should draw and which ones
		// are behind others we have already drawn. Initialize all values in the
		// buffer to be as far back as possible (-math.MaxFloat64)
		zbuff := make([][]float64, w)
		for i := range zbuff {
			zbuff[i] = make([]float64, h)
			for j := range zbuff[i] {
				zbuff[i][j] = -math.MaxFloat64
			}
		}
		// The center, or origin, is at 0,0,0
		ctr := Vertex{0.0, 0.0, 0.0}
		// Which way up is, or pointing in the y direction
		up := Vertex{0.0, 1.0, 0.0}
		// Where we're looking from
		eye := Vertex{math.Sin(mouseXt), math.Sin(mouseYt), math.Cos(mouseXt)}
		// (More documentation needed here)
		z := eye.Sub(ctr).Unit()
		x := up.Cross(z).Unit()
		y := z.Cross(x)

		// For each triangle, draw it
		for i := 0; i < len(r.tv); i++ {
			// Obtain the normal and triangle values
			// from our view
			// (More documentation needed here)
			nrm := r.tn[i].ViewNrm(x, y, z)
			tri := r.tv[i].ViewTri(x, y, z, eye)
			tex := r.tt[i]
			per := tri.Perspective()
			vew := per.Viewport(floatgeom.Point2{float64(w), float64(h)})
			// Actually draw the triangle given the values we've calculated
			TDraw(r.Sprite.GetRGBA(), zbuff, vew, nrm, tex, r.textureData)
		}
	}
	r.lastmouse = mouse.LastEvent
	// Instead of handling the drawing ourselves, let the embedded Sprite which
	// we've populated the color buffer of draw itself.
	r.Sprite.DrawOffset(buff, xOff, yOff)
}
