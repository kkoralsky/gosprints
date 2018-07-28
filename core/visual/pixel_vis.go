package visual

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gobuffalo/packr"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"time"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	log "github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
)

const (
	FONT_PATH = "m50.ttf"
)

var (
	fontSize            float64 = 20
	fontScale           float64 = 10
	fontScaleMax        float64 = 350
	playerNameFontScale float64 = 1.2
	resultsFontScale    float64 = 1
	fontColor                   = colornames.Skyblue
	backgroundColor             = colornames.Black
)

type pixelBaseVis struct {
	BaseVis

	winCfg                *pixelgl.WindowConfig
	win                   *pixelgl.Window
	imd                   *imdraw.IMDraw
	updateRaceFunction    func(playerNum, dist uint32)
	drawDashboardFunction func(playerNum uint32)
	fontAtlas             *text.Atlas
}

func (b *pixelBaseVis) Run() {
	var (
		err  error
		face font.Face
	)
	face, err = loadTTF(FONT_PATH, fontSize)
	if err != nil {
		panic(err)
	}
	b.fontAtlas = text.NewAtlas(face, text.ASCII, text.RangeTable(unicode.Latin))

	pixelgl.Run(func() {
		if b.visCfg.Fullscreen {
			b.winCfg.Monitor = pixelgl.PrimaryMonitor()
		}
		b.win, err = pixelgl.NewWindow(*b.winCfg)
		if err != nil {
			panic(err)
		}

		for !b.win.Closed() {
			time.Sleep(1 * time.Second)
			b.win.UpdateInput()
		}
	})
}

func (b *pixelBaseVis) NewTournament(_ context.Context, tournament *pb.Tournament) (*pb.Empty, error) {
	var (
		color_rgba color.RGBA
		ok         bool
		i          uint32
	)
	b.playerCount = tournament.PlayerCount
	b.destValue = tournament.DestValue
	b.mode = tournament.Mode
	switch b.mode {
	case pb.Tournament_DISTANCE:
		b.modeUnit = "s"
	case pb.Tournament_TIME:
		b.modeUnit = "m"
	default:
		b.modeUnit = ""
		log.ErrorLogger.Printf("couldnt set unit for mode: %v", b.mode)
	}

	if b.playerCount > uint32(len(tournament.Color)) {
		return nil, fmt.Errorf("Not enough color defined for players")
	}

	b.colors = nil
	for i = 0; i < b.playerCount; i++ {
		color_rgba, ok = colornames.Map[tournament.Color[i]]
		if !ok {
			return nil, fmt.Errorf("Color %s is unknown", tournament.Color[i])
		}

		b.colors = append(b.colors, color_rgba)
	}
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) NewRace(_ context.Context, race *pb.Race) (*pb.Empty, error) {
	var winCenter = b.win.Bounds().Center()

	b.Clear()
	b.playerNames = nil
	starterText := text.New(winCenter, b.fontAtlas)
	starterText.LineHeight = b.fontAtlas.LineHeight() * 2.5
	for i, p := range race.Players {
		starterText.Color = b.colors[i]
		starterText.Dot.X -= starterText.BoundsOf(p.Name).W() / 2
		starterText.WriteString(p.Name + "\n")

		b.playerNames = append(b.playerNames, p.Name)
		if i+1 != len(race.Players) {
			starterText.Color = fontColor
			starterText.Dot.X -= starterText.BoundsOf("vs").W() / 2
			starterText.WriteString("vs\n")
		}
	}
	starterText.Draw(b.win, pixel.IM.Scaled(starterText.Bounds().Center(), playerNameFontScale).
		Moved(pixel.V(0, starterText.Bounds().H()/2)))
	b.win.Update()
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) AbortRace(_ context.Context, abortMessage *pb.AbortMessage) (*pb.Empty, error) {
	var (
		winCenter = b.win.Bounds().Center()
	)
	messageText := text.New(winCenter, b.fontAtlas)
	messageText.Color = fontColor

	if abortMessage.Message == "" {
		abortMessage.Message = "aborted"
	}
	messageText.WriteString(abortMessage.Message)
	messageText.Draw(b.win, pixel.IM.Moved(pixel.V(-messageText.Bounds().W()/2,
		-messageText.Bounds().H()/3)).Scaled(winCenter, fontScale))

	b.win.Update()
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) StartRace(_ context.Context, starter *pb.Starter) (*pb.Empty, error) {
	if len(b.playerNames) != int(b.playerCount) {
		return &pb.Empty{}, errors.New("player names not set properly - run NewRace first")
	}
	var (
		frameSleep = time.Duration(1000*starter.CountdownTime/3) * time.Millisecond
		winCenter  = b.win.Bounds().Center()
	)
	countdownText := text.New(winCenter, b.fontAtlas)
	countdownText.Color = fontColor
	countdownText.Orig.X -= countdownText.BoundsOf("3").W() / 2
	countdownText.Orig.Y -= countdownText.BoundsOf("3").H() / 3

	for i := 3; i >= 0; i-- {
		b.drawDashboards()
		countdownText.Clear()
		fmt.Fprintf(countdownText, "%d", i)
		countdownText.Draw(b.win, pixel.IM.Scaled(winCenter, fontScale))
		time.Sleep(frameSleep)
		b.win.Update()
	}

	b.drawDashboards()
	b.scaleGo(fontScale)

	b.win.Update()
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) FinishRace(_ context.Context, results *pb.Results) (*pb.Empty, error) {
	b.Clear()

	if len(results.Result) > len(b.colors) {
		return &pb.Empty{}, fmt.Errorf("not enough colors defined to show all results")
	}
	var (
		winCenter  = b.win.Bounds().Center()
		resultText = text.New(winCenter, b.fontAtlas)
	)
	for i, result := range results.Result {
		resultText.Color = b.colors[i]
		resultText.WriteString(result.Player.Name)
		resultText.Color = fontColor
		fmt.Fprintf(resultText, " %.3f%s\n\n", b.getResult(result.Result), b.modeUnit)
	}

	resultText.Draw(b.win, pixel.IM.Moved(winCenter.Sub(resultText.Bounds().Center())).
		Scaled(winCenter, playerNameFontScale))
	b.win.Update()
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) ShowResults(_ context.Context, results *pb.Results) (*pb.Empty, error) {
	var (
		winCenter = b.win.Bounds().Center()
	)
	b.Clear()

	resultsText := text.New(winCenter, b.fontAtlas)
	resultsText.TabWidth = 50

	for i, result := range results.Result {
		fmt.Fprintf(resultsText, "%3d.%s\t\t%10.3f\n", i+1, result.Player.Name, result.Result)
	}
	resultsText.Draw(b.win, pixel.IM.Moved(winCenter.Sub(resultsText.Bounds().Center())).Scaled(winCenter, resultsFontScale))

	b.win.Update()
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) ConfigureVis(_ context.Context, visCfg *pb.VisConfiguration) (*pb.Empty, error) {
	b.visCfg = visCfg
	if b.ResetConfiguration() {
		b.winCfg = &pixelgl.WindowConfig{
			Title:     fmt.Sprintf("gosprints %s %s", b.visCfg.VisName, visCfg.HostName),
			Bounds:    pixel.R(0, 0, float64(visCfg.ResolutionWidth), float64(visCfg.ResolutionHeight)),
			Resizable: false,
			VSync:     true,
		}
		b.imd = imdraw.New(nil)
	}
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) UpdateRace(stream pb.Visual_UpdateRaceServer) error {
	if len(b.playerNames) != int(b.playerCount) {
		return errors.New("player names not set properly - run NewRace first")
	}

	var (
		racer *pb.Racer
		err   error
	)
	for i := float64(fontScale); ; i++ {
		racer, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}
		if time.Now().Sub(b.racingData[racer.PlayerNum].ts) > time.Second {
			b.updRacingData(racer.PlayerNum, racer.Distance)
		}
		b.updateRaceFunction(racer.PlayerNum, racer.Distance)
		b.win.SetColorMask(colornames.White)
		// if i < fontScaleMax {
		// b.scaleGo(i * 2)
		// }
		b.win.Update()
	}
	return nil
}

func (p *pixelBaseVis) StopVis(context.Context, *pb.Empty) (*pb.Empty, error) {
	p.win.SetClosed(true)
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) scaleGo(scale float64) {
	var (
		winCenter = b.win.Bounds().Center()
	)

	goText := text.New(winCenter, b.fontAtlas)
	goText.Color = fontColor
	goText.WriteString("GO")
	goText.DrawColorMask(b.win, pixel.IM.Moved(pixel.V(-goText.Bounds().W()/2,
		-goText.Bounds().H()/3)).Scaled(winCenter, scale),
		pixel.Alpha(1-(scale-fontScale)/fontScaleMax))
}

func (b *pixelBaseVis) drawDashboards() {
	b.Clear()
	for i := 0; i < int(b.playerCount); i++ {
		b.drawDashboardFunction(uint32(i))
	}
	b.win.SetColorMask(colornames.White)
}

func (b *pixelBaseVis) Clear() {
	b.win.Clear(backgroundColor)
}

func loadPicture(path string) (pixel.Picture, error) {
	box := packr.NewBox("./assets")
	image_bytes, err := box.MustBytes(path)
	if err != nil {
		return nil, err
	}
	log.DebugLogger.Printf("file: %s loaded; %d bytes read", path, len(image_bytes))
	image_reader := bytes.NewReader(image_bytes)
	img, _, err := image.Decode(image_reader)
	if err != nil {
		return nil, err
	}
	log.DebugLogger.Printf("file: %s decoded; size: %d x %d", path, img.Bounds().Dx(), img.Bounds().Dy())
	return pixel.PictureDataFromImage(img), nil
}

func loadSprite(path string) (*pixel.Sprite, error) {
	picture, err := loadPicture(path)
	if err != nil {
		return &pixel.Sprite{}, err
	}
	return pixel.NewSprite(picture, picture.Bounds()), nil
}

func loadTTF(path string, size float64) (font.Face, error) {
	box := packr.NewBox("./assets")
	font_bytes, err := box.MustBytes(path)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(font_bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
