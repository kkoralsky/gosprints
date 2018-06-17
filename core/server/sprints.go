package server

import (
	"context"
	"github.com/kkoralsky/gosprints/core/device"
	pb "github.com/kkoralsky/gosprints/proto"
)

type Sprints struct {
	inputDevice *device.InputDevice
	visClient   *VisClient
	tournament  *pb.Tournament
	curRace     *pb.Race
}

func SetupSprints(device *device.InputDevice, visClient *VisClient) *Sprints {
	return &Sprints{
		inputDevice: device,
		visClient:   visClient,
	}
}

func (s *Sprints) NewTournament(ctx context.Context, tournament *pb.Tournament) (*pb.Tournament, error) {
	s.tournament = tournament
	return tournament, nil
}

func (s *Sprints) NewRace(ctx context.Context, race *pb.Race) (*pb.Empty, error) {
	s.curRace = race
	return &pb.Empty{}, nil
}

func (s *Sprints) StartRace(context.Context, *pb.Empty) (*pb.Empty, error) {
	s.visClient.StartRace()

	return &pb.Empty{}, nil
}

func (s *Sprints) AbortRace(context.Context, *pb.Empty) (*pb.Empty, error) {
	s.visClient.AbortRace()

	return &pb.Empty{}, nil
}

func (s *Sprints) ConfigureVis(context.Context, visCfg *pb.VisConfiguration) (*pb.Empty, error) {
	s.visClient.ConfigureVis(visCfg)

	return &pb.Empty{}, nil
}
