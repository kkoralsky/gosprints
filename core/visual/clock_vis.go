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
		winWidth        = c.win.Bounds().W()
		spriteWidth     = c.clockSprite.Picture().Bounds().W()
		scale           = (winWidth/float64(c.playerCount) - 10) / spriteWidth
		clockWidth      = spriteWidth * scale
		color           = c.colors[int(playerNum)]
		verticalPos     = 5 + float64(playerNum)*clockWidth + clockWidth/2
		horizontalPos   = float64(300)
		racingData, ok  = c.racingData[playerNum]
		playerName      = c.playerNames[playerNum]
		dataText        = text.New(pixel.V(verticalPos, 120), c.fontAtlas)
		playerNameWidth = dataText.BoundsOf(playerName).W()
	)

	c.win.SetColorMask(color)
	dataText.Dot.X -= playerNameWidth / 2
	dataText.WriteString(playerName + "\n\n")
	if ok {
		dataText.Dot.X -= playerNameWidth * 3 / 4
		fmt.Fprintf(dataText, "D:%.2fm\n", racingData.realDist)
		dataText.Dot.X -= playerNameWidth * 3 / 4
		fmt.Fprintf(dataText, "V:%.2fkm/h", racingData.velo)
	}

	c.clockSprite.Draw(c.win, pixel.IM.Scaled(pixel.V(0, 0), scale).Moved(pixel.V(verticalPos, horizontalPos)))
	dataText.Draw(c.win, pixel.IM.Scaled(pixel.V(verticalPos, 120), 2))
}

func (c *clockVis) clearDashboard(playerNum uint32) {
	var (
		winWidth    = c.win.Bounds().W()
		spriteWidth = c.clockSprite.Picture().Bounds().W()
		scale       = winWidth / float64(c.playerCount) / spriteWidth
		clockWidth  = spriteWidth * scale
		verticalMin = float64(playerNum) * clockWidth
		verticalMax = verticalMin + clockWidth
	)
	imd := imdraw.New(nil)
	imd.Color = backgroundColor
	imd.Push(pixel.V(verticalMin, 0))
	imd.Push(pixel.V(verticalMax, c.win.Bounds().Max.Y))
	imd.Rectangle(0)
	imd.Draw(c.win)
}

func (c *clockVis) updateRace(playerNum, dist uint32) {
	var (
		angle            = -2 * math.Pi * float64(dist*c.visCfg.MovingUnit) / 360
		winWidth         = c.win.Bounds().W()
		clockSpriteWidth = c.clockSprite.Picture().Bounds().W()
		scale            = (winWidth/float64(c.playerCount) - 10) / clockSpriteWidth
		clockWidth       = clockSpriteWidth * scale
		verticalPos      = 5 + float64(playerNum)*clockWidth + clockWidth/2
		pos              = pixel.V(verticalPos, 300+38*scale)
	)

	c.clearDashboard(playerNum)
	c.drawDashboard(playerNum)
	c.pointerSprite.Draw(c.win, pixel.IM.Scaled(pixel.ZV, scale).Moved(pos).Rotated(pixel.V(verticalPos, 300-18*scale), angle))
}
