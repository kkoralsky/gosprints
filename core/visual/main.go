package visual

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	//"github.com/kkoralsky/gosprints/core"
)

// VisualCfg is basic visual configuration backed up by pixel lib.
type VisualCfg struct {
	name             string
	resolutionWidth  uint16
	resolutionHeight uint16
	winCfg           *pixelgl.WindowConfig
	fullscreen       bool
}

func (v *VisualCfg) Run() {
	pixelgl.Run(func() {
		if v.fullscreen {
			v.winCfg.Monitor = pixelgl.PrimaryMonitor()
		}
		win, err := pixelgl.NewWindow(*v.winCfg)
		if err != nil {
			panic(err)
		}

		for !win.Closed() {
			win.Update()
		}
	})
}

func SetupVis(name string, movingUnit uint, fullscreen bool) (VisualCfg, error) {
	visCfg := VisualCfg{
		winCfg: &pixelgl.WindowConfig{
			Title:     name,
			Bounds:    pixel.R(0, 0, 640, 480),
			Resizable: false,
			VSync:     true,
		},
		fullscreen: fullscreen,
	}
	return visCfg, nil
}
