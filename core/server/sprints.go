package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
	"github.com/kkoralsky/gosprints/core/device"
	pb "github.com/kkoralsky/gosprints/proto"
	"sort"
	"time"
)

type Sprints struct {
	inputDevice device.InputDevice
	visMux      *VisMux
	tournament  *pb.Tournament
	curRace     *pb.Race
	results     map[pb.Gender][]*pb.Result
	abortRace   chan struct{}
	sprintsDb   *SprintsDb
}

func SetupSprints(device device.InputDevice, visMux *VisMux, sprintsDb *SprintsDb) (s *Sprints) {
	var (
		err        error
		tournament *pb.Tournament
	)
	s = &Sprints{
		inputDevice: device,
		visMux:      visMux,
		sprintsDb:   sprintsDb,
		results:     make(map[pb.Gender][]*pb.Result, 3),
	}
	tournament, err = s.sprintsDb.GetLastTournament()
	if err != nil {
		s.NewTournament(context.Background(), &core.DefultTournament)
	} else {
		s.loadTournament(tournament)
	}

	return s
}

func (s *Sprints) NewTournament(ctx context.Context, tournament *pb.Tournament) (*pb.Tournament, error) {
	s.tournament = tournament
	s.results[pb.Gender_MALE] = []*pb.Result{}
	s.results[pb.Gender_FEMALE] = []*pb.Result{}
	s.results[pb.Gender_OTHER] = []*pb.Result{}

	s.sprintsDb.SaveTournament(s.tournament)

	return tournament, s.visMux.NewTournament(tournament)
}

func (s *Sprints) LoadTournament(ctx context.Context, tournamentSpec *pb.TournamentSpec) (*pb.Tournament, error) {
	tournament, err := s.sprintsDb.GetTournament(tournamentSpec.Name)
	if err != nil {
		return nil, err
	}
	s.loadTournament(tournament)
	return tournament, nil
}

func (s *Sprints) loadTournament(tournament *pb.Tournament) {
	s.tournament = tournament
	s.results[pb.Gender_MALE] = make([]*pb.Result, 0)
	s.results[pb.Gender_FEMALE] = make([]*pb.Result, 0)
	s.results[pb.Gender_OTHER] = make([]*pb.Result, 0)

	for _, result := range s.tournament.Result {
		s.results[result.Player.Gender] = append(s.results[result.Player.Gender], result)
		core.DebugLogger.Printf(
			"%s (%s): %.3f loaded", result.Player.Name,
			pb.Gender_name[int32(result.Player.Gender)],
			result.Result,
		)
	}
	s.visMux.NewTournament(tournament)
}

func (s *Sprints) GetTournamentNames(context.Context, *pb.Empty) (*pb.TournamentNames, error) {
	var tournamentNames = &pb.TournamentNames{Name: []string{}}

	if s.sprintsDb != nil && s.sprintsDb.tournaments != nil {
		for _, tournament := range s.sprintsDb.tournaments.Tournament {
			tournamentNames.Name = append(tournamentNames.Name, tournament.Name)
		}
		return tournamentNames, nil
	}
	return nil, errors.New("No tournaments loaded")
}

func (s *Sprints) GetCurrentTournament(context.Context, *pb.Empty) (*pb.Tournament, error) {
	if s.tournament != nil {
		return s.tournament, nil
	}
	return nil, errors.New("No tournament loaded")

}

func (s *Sprints) ShowResults(_ context.Context, resultSpec *pb.ResultSpec) (*pb.Empty, error) {
	if s.tournament != nil {
		s.visMux.ShowResults(&pb.Results{Result: s.results[resultSpec.Gender]})
		return nil, nil
	}
	return nil, errors.New("No tournament loaded")
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
	var results = make([]*pb.Result, len(s.results[resultSpec.Gender]))
	copy(results, s.results[resultSpec.Gender])
	sort.Slice(results, func(i, j int) bool {
		if s.tournament.Mode == pb.Tournament_TIME {
			return results[i].Result < results[j].Result
		} else {
			return results[i].Result > results[j].Result
		}
	})

	core.DebugLogger.Printf("sending %d sorted results", len(results))

	for _, result := range results {
		if err := stream.Send(result); err != nil {
			return err
		}
	}
	return nil
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
								playersTimes[i] = uint(time.Now().Sub(start) / time.Millisecond)
								core.DebugLogger.Printf("player #%d finished", i)
							}
						}
					}
				}
				time.Sleep(time.Millisecond)
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
						return nil
					default:
						getDistance(i)
					}
				}
				time.Sleep(time.Millisecond)
			}
			return playersDists
		}
		finishRace = func(results map[int]uint) {
			var protoResults []*pb.Result

			for playerNum, result := range results {
				playerGender := s.curRace.Players[playerNum].Gender
				resultPb := &pb.Result{
					DestValue: s.curRace.DestValue,
					Player:    s.curRace.Players[playerNum],
					Result:    float32(result),
				}
				protoResults = append(protoResults, resultPb)
				s.results[playerGender] = append(s.results[playerGender], resultPb)
				s.persistResult(resultPb)
			}
			s.curRace = nil
			s.visMux.FinishRace(&pb.Results{Result: protoResults})
		}
	)

	s.visMux.SetupRacers()

	if s.tournament.Mode == pb.Tournament_TIME {
		results = doTimedRace()
	} else if s.tournament.Mode == pb.Tournament_DISTANCE {
		results = doDistanceRace()
	}

	s.visMux.CloseRacers()
	finishRace(results)
}

func (s *Sprints) persistResult(resultPb *pb.Result) {
	s.tournament.Result = append(s.tournament.Result, resultPb)
	if err := s.sprintsDb.SaveTournament(s.tournament); err != nil {
		core.ErrorLogger.Fatalf("error while saving tournament: %v", err)
	}
}
