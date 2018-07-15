package visual

import (
	"context"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
)

type VisInterface interface {
	pb.VisualServer
	Run()
	ResetConfiguration() bool
	IsConfigured() bool
	GetVisCfg() *pb.VisConfiguration
}

type BaseVis struct {
	VisInterface
	visCfg     *pb.VisConfiguration
	configured bool
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
	resolutionHeight uint, movingUnit uint) (VisInterface, error) {

	vis, err := selectVis(visName)
	if vis != nil {
		vis.ConfigureVis(context.Background(), &pb.VisConfiguration{
			HostName:         hostName,
			VisName:          visName,
			Fullscreen:       fullscreen,
			ResolutionWidth:  uint32(resolutionWidth),
			ResolutionHeight: uint32(resolutionHeight),
			MovingUnit:       uint32(movingUnit),
		})
	}

	return vis, err
}
