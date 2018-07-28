package visual

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	log "github.com/kkoralsky/gosprints/core"
	"math"
)

const (
	STOPWATCH_PATH = "stopwatch.png"
	POINTER_PATH   = "pointer.png"
)

type clockVis struct {
	pixelBaseVis
	clockSprite   *pixel.Sprite
	pointerSprite *pixel.Sprite
	curAngles     map[uint32]float64
}

func NewClockVis() *clockVis {
	var err error
	c := &clockVis{}
	c.clockSprite, err = loadSprite(STOPWATCH_PATH)
	if err != nil {
		log.ErrorLogger.Fatalf("couldnt load sprite: %s", STOPWATCH_PATH)
	}
	c.pointerSprite, err = loadSprite(POINTER_PATH)
	if err != nil {
		log.ErrorLogger.Fatalf("couldnt load sprite: %s", POINTER_PATH)
	}

	c.updateRaceFunction = c.updateRace
	c.drawDashboardFunction = c.drawDashboard
	return c
}

func (c *clockVis) drawDashboard(playerNum uint32) {
	var (
		winWidth      = c.win.Bounds().W()
		spriteWidth   = c.clockSprite.Picture().Bounds().W()
		horizontalPos = winWidth / 2
		scale         = horizontalPos / spriteWidth
		// clockWidth          = spriteWidth * scale
		verticalPos         = float64(300)
		textHorizontalWidth = (winWidth/float64(c.playerCount) - 10)
		textHorizontalPos   = 5 + float64(playerNum)*textHorizontalWidth + textHorizontalWidth/2
		color               = c.colors[int(playerNum)]
		racingData, ok      = c.racingData[playerNum]
		playerName          = c.playerNames[playerNum]
		dataText            = text.New(pixel.V(textHorizontalPos, 120), c.fontAtlas)
		playerNameWidth     = dataText.BoundsOf(playerName).W()
	)

	dataText.Color = color
	dataText.Dot.X -= playerNameWidth / 2
	dataText.WriteString(playerName + "\n\n")
	if ok {
		dataText.Dot.X -= playerNameWidth * 3 / 4
		fmt.Fprintf(dataText, "D:%d\n", racingData.dist)
		dataText.Dot.X -= playerNameWidth * 3 / 4
		fmt.Fprintf(dataText, "V:%.2fkm/h", racingData.velo)
	}

	c.clockSprite.Draw(c.win, pixel.IM.Scaled(pixel.V(0, 0), scale).Moved(pixel.V(horizontalPos, verticalPos)))
	dataText.Draw(c.win, pixel.IM.Scaled(pixel.V(verticalPos, 120), playerNameFontScale))
}

func (c *clockVis) clearDashboard(playerNum uint32) {
	var (
		winWidth            = c.win.Bounds().W()
		winHeight           = c.win.Bounds().H()
		spriteWidth         = c.clockSprite.Picture().Bounds().W()
		spriteHeight        = c.clockSprite.Picture().Bounds().H()
		scale               = winWidth / 2 / spriteWidth
		textHorizontalWidth = (winWidth/float64(c.playerCount) - 10)
		textHorizontalPos   = 5 + float64(playerNum)*textHorizontalWidth
	)
	imd := imdraw.New(nil)
	imd.Color = backgroundColor
	imd.Push(
		pixel.V(0, winHeight),
		pixel.V(winWidth, winHeight-scale*spriteHeight),
	)
	imd.Rectangle(0)
	imd.Push(
		pixel.V(textHorizontalPos, 120),
		pixel.V(textHorizontalPos+textHorizontalWidth, 0),
	)
	imd.Rectangle(0)
	imd.Draw(c.win)
}

func (c *clockVis) updateRace(playerNum, dist uint32) {
	var (
		angle            = -2 * math.Pi * float64(dist*c.visCfg.MovingUnit) / 360
		winWidth         = c.win.Bounds().W()
		clockSpriteWidth = c.clockSprite.Picture().Bounds().W()
		horizontalPos    = winWidth / 2
		scale            = horizontalPos / clockSpriteWidth
		pos              = pixel.V(horizontalPos, 300+38*scale)
		i                uint32
	)

	c.clearDashboard(playerNum)
	c.drawDashboard(playerNum)

	if c.curAngles == nil {
		c.curAngles = make(map[uint32]float64, c.playerCount)
	}
	c.curAngles[playerNum] = angle
	for i, angle = range c.curAngles {
		c.win.SetColorMask(c.colors[int(i)])
		c.pointerSprite.Draw(c.win, pixel.IM.Scaled(pixel.ZV, scale).Moved(pos).Rotated(pixel.V(horizontalPos, 300-18*scale), angle))
	}
}
