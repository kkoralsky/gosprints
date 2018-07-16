package visual

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gobuffalo/packr"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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
	log "github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
)

var (
	fontFace                    = basicfont.Face7x13
	fontScale           float64 = 30
	fontScaleMax        float64 = 350
	playerNameFontScale float64 = 5
	fontColor                   = colornames.Skyblue
	backgroundColor             = colornames.Black
)

type pixelBaseVis struct {
	BaseVis

	winCfg                *pixelgl.WindowConfig
	win                   *pixelgl.Window
	imd                   *imdraw.IMDraw
	updateRaceFunction    func(playerNum, distance uint32)
	drawDashboardFunction func(playerNum uint32)
	colors                []color.RGBA
	playerCount           uint32
	playerNames           []string
	destValue             uint32
	mode                  rune
	fontAtlas             *text.Atlas
	modeUnit              string
}

func (b *pixelBaseVis) Run() {
	b.fontAtlas = text.NewAtlas(fontFace, text.ASCII, text.RangeTable(unicode.Latin))

	b.NewTournament(context.Background(), &pb.Tournament{
		Color:       []string{"blue", "red", "green", "yellow"},
		DestValue:   400,
		Name:        "default tournament",
		Mode:        pb.Tournament_DISTANCE,
		PlayerCount: 2,
	})

	pixelgl.Run(func() {
		var err error
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
	b.mode = rune(pb.Tournament_TournamentMode_name[int32(tournament.Mode)][0])
	switch b.mode {
	case 'd':
		b.modeUnit = "s"
	case 't':
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
	var (
		winWidth    = b.win.Bounds().W()
		winHeight   = b.win.Bounds().H()
		playerSpace = winHeight / float64(len(race.Players))
		lineHeight  = b.fontAtlas.LineHeight()
	)

	b.win.Clear(backgroundColor)
	b.playerNames = nil
	for i, p := range race.Players {
		playerWriter := text.New(pixel.V(winWidth/2, winHeight-float64(i)*playerSpace-playerSpace/2+lineHeight/2), b.fontAtlas)
		playerWriter.Color = b.colors[i]
		playerWriter.Orig.X -= playerWriter.BoundsOf(p.Name).W() / 2
		playerWriter.WriteString(p.Name)
		playerWriter.Draw(b.win, pixel.IM.Scaled(playerWriter.Bounds().Center(), playerNameFontScale))
		b.playerNames = append(b.playerNames, p.Name)

		if i+1 != len(race.Players) {
			vsWriter := text.New(pixel.V(winWidth/2-10, winHeight-float64(i+1)*playerSpace+lineHeight/2), b.fontAtlas)
			vsWriter.WriteString("vs")
			vsWriter.Draw(b.win, pixel.IM.Scaled(vsWriter.Bounds().Center(), playerNameFontScale))
		}
	}
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
		-messageText.Bounds().H()/3)).Scaled(winCenter, 7))

	b.win.Update()
	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) StartRace(_ context.Context, starter *pb.Starter) (*pb.Empty, error) {
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
	if len(results.Result) > len(b.colors) {
		return &pb.Empty{}, fmt.Errorf("not enough colors defined to show all results")
	}
	for i, result := range results.Result {
		h := float64(b.visCfg.ResolutionHeight - (uint32(i)+1)*20)
		playerText := text.New(pixel.V(float64(b.visCfg.ResolutionWidth/3), h), b.fontAtlas)
		playerText.Color = b.colors[i]
		playerText.WriteString(result.Player.Name)
		playerText.Draw(b.win, pixel.IM)

		resultText := text.New(pixel.V(float64(b.visCfg.ResolutionWidth/2), h), b.fontAtlas)
		fmt.Fprintf(resultText, "%.3f%s", result.Result, b.modeUnit)
		resultText.Draw(b.win, pixel.IM)
	}

	return &pb.Empty{}, nil
}

func (b *pixelBaseVis) ShowResults(_ context.Context, results *pb.Results) (*pb.Empty, error) {
	var (
		winCenter = b.win.Bounds().Center()
	)
	b.win.Clear(backgroundColor)

	resultsText := text.New(winCenter, b.fontAtlas)
	resultsText.TabWidth = 50

	for i, result := range results.Result {
		fmt.Fprintf(resultsText, "%3d.\t%s\t\t\t\t\t\t\t\t%10.3f\n", i+1, result.Player.Name, result.Result)
	}
	resultsText.Draw(b.win, pixel.IM.Moved(winCenter.Sub(resultsText.Bounds().Center())).Scaled(winCenter, 1.5))

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
		b.updateRaceFunction(racer.PlayerNum, racer.Distance)
		b.win.SetColorMask(colornames.White)
		if i < fontScaleMax {
			b.scaleGo(i)
		}
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
	b.win.Clear(backgroundColor)
	for i := 0; i < int(b.playerCount); i++ {
		b.drawDashboardFunction(uint32(i))
	}
	b.win.SetColorMask(colornames.White)
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
