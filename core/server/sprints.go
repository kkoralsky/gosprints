package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
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
	abortRace   chan struct{}
}

func SetupSprints(device device.InputDevice, visMux *VisMux) *Sprints {
	return &Sprints{
		inputDevice: device,
		visMux:      visMux,
		tournament:  &core.DefultTournament,
	}
}

func (s *Sprints) NewTournament(ctx context.Context, tournament *pb.Tournament) (*pb.Tournament, error) {
	s.tournament = tournament
	core.DebugLogger.Printf("%s", pb.Tournament_TournamentMode_name[int32(s.tournament.Mode)])
	return tournament, s.visMux.NewTournament(tournament)
}

func (s *Sprints) NewRace(ctx context.Context, race *pb.Race) (*pb.Empty, error) {
	s.curRace = race
	err := s.visMux.NewRace(race)
	return &pb.Empty{}, err
}

func (s *Sprints) StartRace(_ context.Context, starter *pb.Starter) (*pb.Player, error) {
	if s.curRace == nil {
		return &pb.Player{}, errors.New("race is not established")
	}
	s.visMux.StartRace(starter)
	s.inputDevice.Clean()
	s.abortRace = make(chan struct{}, 10) // allow to up to 10 unhandled "aborts"
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
	if s.curRace != nil {
		s.abortRace <- struct{}{}
		s.visMux.AbortRace(abortMessage)
	}

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
	var (
		playersCount = len(s.curRace.Players)
		playersDists = make(map[int]uint, playersCount)
		results      map[int]uint
		getDistance  = func(i int) uint {
			if dist, err := s.inputDevice.GetDist(uint(i)); err != nil {
				core.DebugLogger.Printf(err.Error())
			} else if dist != playersDists[i] {
				playersDists[i] = dist
				s.visMux.SendRaceUpdate(uint32(i), uint32(dist))
			}
			return playersDists[i]
		}
		doDistanceRace = func() (playersTimes map[int]uint) {
			var (
				wholeDistance = uint(s.curRace.DestValue)
				start         = time.Now()
			)
			playersTimes = make(map[int]uint, playersCount)

			for playersFinished := 0; playersFinished < playersCount; {
				for i := 0; i < playersCount; i++ {
					select {
					case _, ok := <-s.abortRace:
						if ok {
							close(s.abortRace)
						}
						return
					default:
						if getDistance(i) >= wholeDistance {
							if playersTimes[i] == 0 {
								playersFinished++
								playersTimes[i] = uint(time.Now().Sub(start))
								core.DebugLogger.Printf("player #%d finished", i)
							}
						}
					}
				}
				time.Sleep(50 * time.Millisecond)
			}
			return
		}
		doTimedRace = func() map[int]uint {
			var (
				start        = time.Now()
				raceDuration = time.Duration(s.curRace.DestValue) * time.Second
				finish       = start.Add(raceDuration)
			)

			for now := time.Now(); now.Before(finish); now = time.Now() {
				for i := 0; i < playersCount; i++ {
					select {
					case <-s.abortRace:
						close(s.abortRace)
						return playersDists
					default:
						getDistance(i)
					}
				}
				time.Sleep(50 * time.Millisecond)
			}
			return playersDists
		}
		finishRace = func(results map[int]uint) {
			var protoResults []*pb.Result

			for playerNum, result := range results {
				protoResults = append(protoResults, &pb.Result{
					DestValue: s.curRace.DestValue,
					Player:    s.curRace.Players[playerNum],
					Result:    float32(result),
				})
			}

			s.curRace = nil
			s.visMux.FinishRace(&pb.Results{Result: protoResults})
		}
	)

	s.visMux.SetupRacers()

	core.DebugLogger.Printf("%s", pb.Tournament_TournamentMode_name[int32(s.tournament.Mode)])
	if s.tournament.Mode == pb.Tournament_TIME {
		results = doTimedRace()
	} else if s.tournament.Mode == pb.Tournament_DISTANCE {
		results = doDistanceRace()
	}

	s.visMux.CloseRacers()
	finishRace(results)
}
