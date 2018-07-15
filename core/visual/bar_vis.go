package visual

import (
	// "github.com/faiface/pixel"
	// "github.com/faiface/pixel/imdraw"
	// "github.com/faiface/pixel/pixelgl"
	"github.com/kkoralsky/gosprints/core"
)

type barVis struct {
	pixelBaseVis
}

func NewBarVis() *barVis {
	b := &barVis{}
	b.updateRaceFunction = b.updateRace
	b.drawDashboardFunction = b.drawDashboard
	return b
}

func (b *barVis) drawDashboard(playerNum uint32) {

}

func (b *barVis) updateRace(playerNum, distance uint32) {
	core.DebugLogger.Printf("updating player %d with distance: %d",
		playerNum, distance)
}
