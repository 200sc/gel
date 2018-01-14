package gel

import "image"

type Target struct {
	vew, nrm, tex Triangle
	fdif          *image.RGBA
}
