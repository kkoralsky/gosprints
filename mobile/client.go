package main

import (
	"context"
	"fmt"
	log "github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"github.com/therecipe/qt/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"io"

	"time"
)

// SprintsClient is GRPC client calling Sprints service
type SprintsClient struct {
	core.QObject

	conn        *grpc.ClientConn
	addr        string
	connState   connectivity.State
	client      pb.SprintsClient
	tournament  *pb.Tournament
	race        *pb.Race
	resultModel *ResultModel

	_ func(msg string)      `signal:"info"`
	_ func(err, msg string) `signal:"error"`
	_ func(msg string)      `signal:"success"`

	_ int                                                 `property:"connState"`
	_ func(string, uint, bool) string                     `slot:"dialGrpc"`
	_ func(string, uint, int32, uint, []string) string    `slot:"newTournament"`
	_ func([]string, uint) string                         `slot:"newRace"`
	_ func() string                                       `slot:"startRace"`
	_ func() string                                       `slot:"abortRace"`
	_ func(string, string, bool, uint, uint, uint) string `slot:"configureVis"`
	_ func(string)                                        `slot:"getResults"`
}

type SprintsClientInterface interface {
	dialGrpc(string, uint, bool) string
	newTournament(string, uint, int32, uint, []string) string
	newRace([]string, uint) string
	startRace() string
	abortRace() string
	configureVis(string, string, bool, uint, uint, uint) string
	getResults(string)
}

func SetupSprintsClient(resultModel *ResultModel) *SprintsClient {
	client := NewSprintsClient(nil)

	client.resultModel = resultModel
	client.connState = connectivity.Shutdown

	client.ConnectDialGrpc(client.dialGrpc)
	client.ConnectNewTournament(client.newTournament)
	client.ConnectNewRace(client.newRace)
	client.ConnectStartRace(client.startRace)
	client.ConnectAbortRace(client.abortRace)
	client.ConnectConfigureVis(client.configureVis)
	client.ConnectGetResults(client.getResults)

	// client.dialGrpc(defaultHost, defaultPort, false)

	return client
}

func (s *SprintsClient) dialGrpc(hostName string, port uint, blocking bool) string {
	var err error
	s.Close()

	s.addr = fmt.Sprintf("%s:%d", hostName, port)
	log.DebugLogger.Printf("trying connect to %s endpoint\n", s.addr)
	if blocking {
		s.conn, err = grpc.Dial(s.addr, grpc.WithInsecure(), grpc.WithTimeout(10*time.Second),
			grpc.WithBlock())
	} else {
		s.conn, err = grpc.Dial(s.addr, grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
	}
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return err.Error()
	}

	s.client = pb.NewSprintsClient(s.conn)
	s.setConnState()
	log.InfoLogger.Printf("connection %s dialed to %s", s.connState.String(), s.addr)
	go s.updateConnectionState()

	return ""
}

func (s *SprintsClient) setConnState() {
	s.connState = s.conn.GetState()
	s.SetConnState(int(s.connState)) // notify qml
}

func (s *SprintsClient) updateConnectionState() {
	for s.conn != nil {
		if s.conn.WaitForStateChange(context.Background(), s.connState) {
			if s.conn != nil {
				s.setConnState()
			} else {
				break
			}
			log.DebugLogger.Printf("connection state changed: %s", s.connState.String())
		} else {
			log.InfoLogger.Println("connection expired")
			break
		}
	}
	log.DebugLogger.Printf("stop updating connection state: %s", s.connState.String())
}

func (s *SprintsClient) newTournament(name string, destValue uint, mode int, playerCount uint, colors []string) string {
	var err error
	s.tournament, err = s.client.NewTournament(context.Background(), &pb.Tournament{
		Name:      name,
		DestValue: uint32(destValue),
		// Mode:        pb.Tournament_TournamentMode(pb.Tournament_TournamentMode_value[string(mode)]),
		Mode:        pb.Tournament_TournamentMode(mode),
		Color:       colors,
		PlayerCount: uint32(playerCount),
	})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) newRace(playerNames []string, destValue uint) string {
	var players []*pb.Player
	for _, playerName := range playerNames {
		players = append(players, &pb.Player{
			Name: playerName,
		})
	}
	_, err := s.client.NewRace(context.Background(), &pb.Race{
		Players:   players,
		DestValue: uint32(destValue),
	})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) startRace() string {
	_, err := s.client.StartRace(context.Background(), &pb.Starter{CountdownTime: 3})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) abortRace() string {
	_, err := s.client.AbortRace(context.Background(), &pb.AbortMessage{Message: "aborted"})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) configureVis(hostName string, visName string, fullscreen bool, resolutionWidth uint, resolutionHeight uint, movingUnit uint) string {
	_, err := s.client.ConfigureVis(context.Background(), &pb.VisConfiguration{
		HostName:         hostName,
		VisName:          visName,
		Fullscreen:       fullscreen,
		ResolutionWidth:  uint32(resolutionWidth),
		ResolutionHeight: uint32(resolutionHeight),
		MovingUnit:       uint32(movingUnit),
	})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) getResults(gender string) {
	var (
		result        *pb.Result
		err           error
		resultsStream pb.Sprints_GetResultsClient
	)
	log.DebugLogger.Printf("getting results for %s", gender)
	resultsStream, err = s.client.GetResults(context.Background(), &pb.ResultSpec{
		Gender:         pb.Gender(pb.Gender_value[gender]),
		Last:           0,
		TournamentName: "",
	})
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return
	}
	s.resultModel.modelReset()
	for {
		result, err = resultsStream.Recv()
		if err != nil {
			if err != io.EOF {
				log.ErrorLogger.Printf("error on results receive: %s", err.Error())
			}
			break
		}
		s.resultModel.AddResult(
			result.Player.Name,
			pb.Gender_name[int32(result.Player.Gender)],
			result.Result,
			uint(result.DestValue))
	}
	return
}

func (s *SprintsClient) Close() bool {
	if s.conn != nil && s.connState != connectivity.Idle && s.connState != connectivity.Ready {
		log.DebugLogger.Printf("closing connection %s", s.addr)
		s.conn.Close()
		s.conn = nil
		return true
	}
	s.conn = nil
	return false
}
