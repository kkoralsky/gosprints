package visual

import (
	"github.com/faiface/pixel/pixelgl"
	pb "github.com/kkoralsky/gosprints/proto"
)

type gameVis struct {
	VisInterface
	visCfg *pb.VisConfiguration
	win    *pixelgl.Window
	winCfg *pixelgl.WindowConfig
}

func NewGameVis() *gameVis {
	return &gameVis{}
}
