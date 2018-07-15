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
	results     map[int32][]*pb.Result
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
	err := s.visClient.NewRace(race)
	return &pb.Empty{}, err
}

func (s *Sprints) StartRace(_ context.Context, starter *pb.Starter) (*pb.Empty, error) {
	s.visClient.StartRace(starter)

	return &pb.Empty{}, nil
}

func (s *Sprints) AbortRace(_ context.Context, abortMessage *pb.AbortMessage) (*pb.Empty, error) {
	s.visClient.AbortRace(abortMessage)

	return &pb.Empty{}, nil
}

func (s *Sprints) ConfigureVis(_ context.Context, visCfg *pb.VisConfiguration) (*pb.Empty, error) {
	s.visClient.ConfigureVis(visCfg)

	return &pb.Empty{}, nil
}

func (s *Sprints) GetResults(resultSpec *pb.ResultSpec, stream pb.Sprints_GetResultsServer) error {
	for _, result := range s.results[int32(resultSpec.Gender)] {
		if err := stream.Send(result); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sprints) GetTournaments(*pb.Empty, pb.Sprints_GetTournamentsServer) error {
	panic("not implemented")
}
