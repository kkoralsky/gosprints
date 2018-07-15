package visual

import (
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

var (
	clockSprite   *pixel.Sprite
	pointerSprite *pixel.Sprite
)

type clockVis struct {
	pixelBaseVis
}

func init() {
	var err error
	clockSprite, err = loadSprite(STOPWATCH_PATH)
	if err != nil {
		log.ErrorLogger.Fatalf("couldnt load sprite: %s", STOPWATCH_PATH)
	}
	pointerSprite, err = loadSprite(POINTER_PATH)
	if err != nil {
		log.ErrorLogger.Fatalf("couldnt load sprite: %s", POINTER_PATH)
	}
}

func NewClockVis() *clockVis {
	c := &clockVis{}
	c.updateRaceFunction = c.updateRace
	c.drawDashboardFunction = c.drawDashboard
	return c
}

func (c *clockVis) drawDashboard(playerNum uint32) {
	var (
		winWidth      = c.win.Bounds().W()
		spriteWidth   = clockSprite.Picture().Bounds().W()
		scale         = (winWidth/float64(c.playerCount) - 10) / spriteWidth
		clockWidth    = spriteWidth * scale
		color         = c.colors[int(playerNum)]
		verticalPos   = 5 + float64(playerNum)*clockWidth + clockWidth/2
		horizontalPos = float64(300)
	)

	c.win.SetColorMask(color)
	playerNameText := text.New(pixel.V(verticalPos, 120), c.fontAtlas)
	playerNameText.Dot.X -= playerNameText.BoundsOf(c.playerNames[playerNum]).W() / 2
	playerNameText.WriteString(c.playerNames[playerNum])
	clockSprite.Draw(c.win, pixel.IM.Scaled(pixel.V(0, 0), scale).Moved(pixel.V(verticalPos, horizontalPos)))
	playerNameText.Draw(c.win, pixel.IM.Scaled(pixel.V(verticalPos, 120), 2))

}

func (c *clockVis) clearDashboard(playerNum uint32) {
	var (
		winWidth    = c.win.Bounds().W()
		spriteWidth = clockSprite.Picture().Bounds().W()
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

func (c *clockVis) updateRace(playerNum, distance uint32) {
	var (
		angle            = -2 * math.Pi * float64(distance*c.visCfg.MovingUnit/10) / 360
		winWidth         = c.win.Bounds().W()
		clockSpriteWidth = clockSprite.Picture().Bounds().W()
		scale            = (winWidth/float64(c.playerCount) - 10) / clockSpriteWidth
		clockWidth       = clockSpriteWidth * scale
		verticalPos      = 5 + float64(playerNum)*clockWidth + clockWidth/2
		pos              = pixel.V(verticalPos, 300+38*scale)
	)

	c.clearDashboard(playerNum)
	c.drawDashboard(playerNum)
	pointerSprite.Draw(c.win, pixel.IM.Scaled(pixel.ZV, scale).Moved(pos).Rotated(pixel.V(verticalPos, 300-18*scale), angle))
}