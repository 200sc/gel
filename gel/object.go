package gel

import (
	"bufio"
	"fmt"
	"os"

	"github.com/oakmound/oak/dlog"
)

type Obj struct {
	vsv, vsn, vst []Vertex
	fs            []Face
}

func (o Obj) Tvgen() []Triangle {
	scale := 1.0 / VMaxLen(o.vsv)
	tv := make([]Triangle, len(o.fs))
	for i, f := range o.fs {
		tv[i] = Triangle{
			o.vsv[f.va],
			o.vsv[f.vb],
			o.vsv[f.vc],
		}.Mul(scale)
	}
	return tv
}

func (o Obj) Tngen() []Triangle {
	tn := make([]Triangle, len(o.fs))
	for i, f := range o.fs {
		tn[i] = Triangle{
			o.vsn[f.na],
			o.vsn[f.nb],
			o.vsn[f.nc],
		}
	}
	return tn
}

func (o Obj) Ttgen() []Triangle {
	tt := make([]Triangle, len(o.fs))
	for i, f := range o.fs {
		tt[i] = Triangle{
			o.vst[f.ta],
			o.vst[f.tb],
			o.vst[f.tc],
		}
	}
	return tt
}

func oparse(f *os.File) Obj {
	vsv := make([]Vertex, 0)
	vsn := make([]Vertex, 0)
	vst := make([]Vertex, 0)
	fs := make([]Face, 0)

	scn := bufio.NewScanner(f)
	defer func() {
		err := f.Close()
		if err != nil {
			dlog.Error(err)
		}
	}()

	for scn.Scan() {
		var f Face
		var v Vertex

		line := scn.Text()

		if len(line) < 2 {
			continue
		}
		if line[0] == 'v' && line[1] == 'n' {
			fmt.Sscanf(line, "vn %f %f %f", &v.x, &v.y, &v.z)
			vsn = append(vsn, v)
		} else if line[0] == 'v' && line[1] == 't' {
			fmt.Sscanf(line, "vt %f %f %f", &v.x, &v.y, &v.z)
			vst = append(vst, v)
		} else if line[0] == 'v' {
			fmt.Sscanf(line, "v %f %f %f", &v.x, &v.y, &v.z)
			vsv = append(vsv, v)
		} else if line[0] == 'f' {
			fmt.Sscanf(line, "f %d/%d/%d %d/%d/%d %d/%d/%d", &f.va, &f.ta, &f.na, &f.vb, &f.tb, &f.nb, &f.vc, &f.tc, &f.nc)
			fs = append(fs, Face{
				f.va - 1, f.vb - 1, f.vc - 1,
				f.ta - 1, f.tb - 1, f.tc - 1,
				f.na - 1, f.nb - 1, f.nc - 1,
			})
		}
	}
	return Obj{vsv, vsn, vst, fs}
}
