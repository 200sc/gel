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
	for i := 0; i < len(o.fs); i++ {
		tv[i] = Triangle{
			o.vsv[o.fs[i].va],
			o.vsv[o.fs[i].vb],
			o.vsv[o.fs[i].vc],
		}.Mul(scale)
	}
	return tv
}

func (o Obj) Tngen() []Triangle {
	scale := 1.0 / VMaxLen(o.vsv)
	tn := make([]Triangle, len(o.fs))
	for i := 0; i < len(o.fs); i++ {
		tn[i] = Triangle{
			o.vsn[o.fs[i].na],
			o.vsn[o.fs[i].nb],
			o.vsn[o.fs[i].nc],
		}.Mul(scale)
	}
	return tn
}

func (o Obj) Ttgen() []Triangle {
	scale := 1.0 / VMaxLen(o.vsv)
	tt := make([]Triangle, len(o.fs))
	for i := 0; i < len(o.fs); i++ {
		tt[i] = Triangle{
			o.vst[o.fs[i].ta],
			o.vst[o.fs[i].tb],
			o.vst[o.fs[i].tc],
		}.Mul(scale)
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
