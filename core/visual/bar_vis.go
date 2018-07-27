package visual

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
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
	var (
		winHeight        = b.win.Bounds().H()
		winWidth         = b.win.Bounds().W()
		color            = b.colors[int(playerNum)]
		playerSpace      = winHeight / float64(b.playerCount)
		playerStartH     = (float64(playerNum) + .5) * playerSpace
		barHeight        = .3 * playerSpace
		imd              = imdraw.New(nil)
		horizontalMargin = .1 * winWidth
		dataText         = text.New(pixel.V(horizontalMargin, playerStartH-2*b.fontAtlas.LineHeight()), b.fontAtlas)
		racingData, ok   = b.racingData[playerNum]
	)

	imd.Color = fontColor
	imd.Push(
		pixel.V(horizontalMargin, playerStartH+barHeight),
		pixel.V(winWidth-horizontalMargin, playerStartH),
	)
	imd.Rectangle(1)
	imd.Draw(b.win)

	dataText.Color = color
	dataText.WriteString(b.playerNames[playerNum])
	dataText.Color = fontColor
	if ok {
		fmt.Fprintf(dataText, "\tD:%.2fm\tV:%.2fkm/h", racingData.realDist, racingData.velo)
	}

	dataText.Draw(b.win, pixel.IM.Scaled(dataText.Bounds().Center(), playerNameFontScale))
	// b.win.SetColorMask(color)
}

func (b *barVis) updateRace(playerNum, distance uint32) {
}
