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

func (b *barVis) clearDashboard(playerNum uint32) {
	var (
		winWidth     = b.win.Bounds().W()
		winHeight    = b.win.Bounds().H()
		playerSpace  = winHeight / float64(b.playerCount)
		playerStartH = (float64(playerNum) + .5) * playerSpace
		imd          = imdraw.New(nil)
	)

	imd.Color = backgroundColor
	imd.Push(
		pixel.V(0, playerStartH),
		pixel.V(winWidth, playerStartH-4*b.fontAtlas.LineHeight()),
	)
	imd.Rectangle(0)
	imd.Draw(b.win)
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
		dataText         = text.New(pixel.V(horizontalMargin,
			playerStartH-2*b.fontAtlas.LineHeight()), b.fontAtlas)
		racingData, ok = b.racingData[playerNum]
	)

	imd.Color = fontColor
	imd.Push(
		pixel.V(horizontalMargin, playerStartH+barHeight),
		pixel.V(winWidth-horizontalMargin, playerStartH),
	)
	imd.Rectangle(1)
	imd.Draw(b.win)

	dataText.Color = color
	// dataText.TabWidth = 3.0
	dataText.WriteString(b.playerNames[playerNum])
	dataText.Color = fontColor
	if ok {
		dataText.Orig.X = winWidth / 2
		dataText.Dot.X = winWidth / 2
		fmt.Fprintf(dataText, "D:%d\nV:%.2fkm/h", racingData.dist, racingData.velo)
	}

	dataText.Draw(b.win, pixel.IM.Scaled(dataText.Bounds().Center(), playerNameFontScale))
	// b.win.SetColorMask(color)
}

func (b *barVis) updateRace(playerNum, distance uint32) {
	var (
		winHeight        = b.win.Bounds().H()
		winWidth         = b.win.Bounds().W()
		horizontalMargin = .1*winWidth + 1
		color            = b.colors[int(playerNum)]
		playerSpace      = winHeight / float64(b.playerCount)
		playerStartH     = (float64(playerNum) + .5) * playerSpace
		barHeight        = .3*playerSpace - 1
		barWidth         = winWidth - 2*horizontalMargin
		curBarWidth      = barWidth * float64(distance) / float64(b.destValue)
		imd              = imdraw.New(nil)
	)

	if curBarWidth > barWidth {
		curBarWidth = barWidth
	}

	b.clearDashboard(playerNum)
	b.drawDashboard(playerNum)
	imd.Color = color
	imd.Push(
		pixel.V(horizontalMargin, playerStartH+barHeight),
		pixel.V(horizontalMargin+curBarWidth, playerStartH),
	)
	imd.Rectangle(0)
	imd.Draw(b.win)
}
