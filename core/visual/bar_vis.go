package visual

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	pb "github.com/kkoralsky/gosprints/proto"
)

type barVis struct {
	//VisInterface
	BaseVis

	winCfg *pixelgl.WindowConfig
	win    *pixelgl.Window
}

func (b *barVis) Run() {
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
			b.win.Update()
		}
	})
}

func (b *barVis) NewTournament(context.Context, *pb.Tournament) (*pb.Tournament, error) {
	panic("not implemented")
}

func (b *barVis) NewRace(context.Context, *pb.Race) (*pb.Empty, error) {
	panic("not implemented")
}

func (b *barVis) StartRace(context.Context, *pb.Empty) (*pb.Empty, error) {
	panic("not implemented")
}

func (b *barVis) AbortRace(context.Context, *pb.Empty) (*pb.Empty, error) {
	panic("not implemented")
}

func (b *barVis) UpdateRace(pb.Races_UpdateRaceServer) error {
	panic("not implemented")
}

func (b *barVis) FinishRace(context.Context, *pb.Empty) (*pb.Empty, error) {
	panic("not implemented")
}

func (b *barVis) ConfigureVis(ctx context.Context, visCfg *pb.VisConfiguration) (*pb.Empty, error) {
	b.visCfg = visCfg
	if b.ResetConfiguration() {
		b.winCfg = &pixelgl.WindowConfig{
			Title:     visCfg.HostName,
			Bounds:    pixel.R(0, 0, float64(visCfg.ResolutionWidth), float64(visCfg.ResolutionHeight)),
			Resizable: false,
			VSync:     true,
		}
	}
	return &pb.Empty{}, nil
}
