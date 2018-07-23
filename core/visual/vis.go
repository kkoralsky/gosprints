package visual

import (
	"context"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"image/color"
	"math"
	"time"
)

type VisInterface interface {
	pb.VisualServer
	Run()
	ResetConfiguration() bool
	IsConfigured() bool
	GetVisCfg() *pb.VisConfiguration
	updRacingData(playerNum, distance uint32) (realDistance, velocity float32)
}

type RacingData struct {
	dist     uint32
	ts       time.Time
	realDist float32
	velo     float32
}

type BaseVis struct {
	VisInterface
	visCfg     *pb.VisConfiguration
	configured bool

	colors      []color.RGBA
	playerCount uint32
	playerNames []string
	destValue   uint32
	mode        pb.Tournament_TournamentMode
	modeUnit    string
	racingData  map[uint32]RacingData
}

func (b *BaseVis) ResetConfiguration() bool {
	b.configured = !b.configured
	return b.configured
}

func (b *BaseVis) IsConfigured() bool {
	return b.configured
}

func (b *BaseVis) GetVisCfg() *pb.VisConfiguration {
	return b.visCfg
}

func (b *BaseVis) updRacingData(playerNum, dist uint32) (realDist, velo float32) {

	var (
		now            = time.Now()
		racingData, ok = b.racingData[playerNum]
		km             = float32(100 * 1000)
		m              = float32(100)
	)
	if ok {
		realDist = float32(b.visCfg.DistFactor*dist) / m
		velo = float32(dist-racingData.dist) / float32(now.Sub(racingData.ts).Hours())
		velo *= float32(b.visCfg.DistFactor) / km
	}

	if b.racingData == nil {
		b.racingData = make(map[uint32]RacingData, b.playerCount)
	}

	b.racingData[playerNum] = RacingData{
		dist:     dist,
		ts:       now,
		realDist: realDist,
		velo:     velo,
	}

	return
}

func (b *BaseVis) getResult(result float32) float32 {
	switch b.mode {
	case pb.Tournament_DISTANCE:
		return result * float32(math.Pow10(-9)) // decode from nanoseconds to seconds
	case pb.Tournament_TIME:
		return float32(b.visCfg.DistFactor) * result / 100 // in meters
	}
	return 0.0
}

func selectVis(visName string) (VisInterface, error) {
	switch visName {
	case "bar":
		return NewBarVis(), nil
	case "clock":
		return NewClockVis(), nil
	case "game":
		return NewGameVis(), nil
	default:
		err := fmt.Errorf("'%s' visualization not found; falling back to 'bar'", visName)
		core.ErrorLogger.Println(err.Error())
		return &barVis{}, nil
	}
}

func Reconfigure(b VisInterface) VisInterface {
	if !b.IsConfigured() {
		vis, _ := selectVis(b.GetVisCfg().VisName)
		vis.ConfigureVis(context.Background(), b.GetVisCfg())
		return vis
	}
	return nil
}

func SetupVis(hostName string, visName string, fullscreen bool, resolutionWidth uint,
	resolutionHeight uint, movingUnit uint, distFactor uint) (VisInterface, error) {

	vis, err := selectVis(visName)
	if vis != nil {
		vis.ConfigureVis(context.Background(), &pb.VisConfiguration{
			HostName:         hostName,
			VisName:          visName,
			Fullscreen:       fullscreen,
			ResolutionWidth:  uint32(resolutionWidth),
			ResolutionHeight: uint32(resolutionHeight),
			MovingUnit:       uint32(movingUnit),
			DistFactor:       uint32(distFactor),
		})

		vis.NewTournament(context.Background(), &core.DefultTournament)
	}

	return vis, err
}
