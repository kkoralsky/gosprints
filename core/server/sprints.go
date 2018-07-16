package server

import (
	"context"
	"fmt"
	log "github.com/kkoralsky/gosprints/core"
	"github.com/kkoralsky/gosprints/core/device"
	pb "github.com/kkoralsky/gosprints/proto"
	"time"
)

type Sprints struct {
	inputDevice device.InputDevice
	visMux      *VisMux
	tournament  *pb.Tournament
	curRace     *pb.Race
	results     map[int32][]*pb.Result
}

func SetupSprints(device device.InputDevice, visMux *VisMux) *Sprints {
	return &Sprints{
		inputDevice: device,
		visMux:      visMux,
	}
}

func (s *Sprints) NewTournament(ctx context.Context, tournament *pb.Tournament) (*pb.Tournament, error) {
	s.tournament = tournament
	return tournament, s.visMux.NewTournament(tournament)
}

func (s *Sprints) NewRace(ctx context.Context, race *pb.Race) (*pb.Empty, error) {
	s.curRace = race
	err := s.visMux.NewRace(race)
	return &pb.Empty{}, err
}

func (s *Sprints) StartRace(_ context.Context, starter *pb.Starter) (*pb.Player, error) {
	s.visMux.StartRace(starter)
	s.inputDevice.Clean()
	time.Sleep(time.Duration(starter.CountdownTime) * time.Second)
	if playerNum, err := s.inputDevice.Check(); err != nil {
		return &pb.Player{}, err
	} else {
		if playerNum >= 0 {
			s.visMux.AbortRace(&pb.AbortMessage{Message: fmt.Sprintf("%s false-started", s.curRace.Players[playerNum].Name)})
			return s.curRace.Players[playerNum], nil
		}
	}

	go s.doRace()

	return &pb.Player{}, nil
}

func (s *Sprints) AbortRace(_ context.Context, abortMessage *pb.AbortMessage) (*pb.Empty, error) {
	s.visMux.AbortRace(abortMessage)

	return &pb.Empty{}, nil
}

func (s *Sprints) ConfigureVis(_ context.Context, visCfg *pb.VisConfiguration) (*pb.Empty, error) {
	s.visMux.ConfigureVis(visCfg)

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

func (s *Sprints) doRace() {
}

func (s *Sprints) doDistanceRace() {
}

func (s *Sprints) doTimedRace() {
	var (
		start        = time.Now()
		raceDuration = time.Duration(s.curRace.DestValue) * time.Second
		finish       = start.Add(raceDuration)
		playerCount  = len(s.curRace.Players)
		dists        []uint
	)

	s.visMux.SetupRacers()

	for i := 0; i < playerCount; i++ {
		dists[i] = 0
	}

	for now := time.Now(); now.Before(finish); now = time.Now() {
		for i := 0; i < len(s.curRace.Players); i++ {
			dist, err := s.inputDevice.GetDist(uint(i))
			if err != nil {
				log.DebugLogger.Printf(err.Error())
				continue
			}
			if dist != dists[i] {
				dists[i] = dist
				s.visMux.SendRaceUpdate(uint32(i), uint32(dist))
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	s.visMux.CloseRacers()
}
