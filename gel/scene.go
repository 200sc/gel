package gel

import (
	"path/filepath"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
)

func gelStart(string, interface{}) {
	r, err := NewRender(filepath.Join("model", "salesman.obj"), "salesman.png")
	if err != nil {
		dlog.Error(err)
		return
	}
	render.Draw(r)
}
func gelLoop() bool {
	return true
}
func gelEnd() (string, *scene.Result) {
	return "gel", nil
}

var Scene = scene.Scene{
	gelStart,
	gelLoop,
	gelEnd,
}
