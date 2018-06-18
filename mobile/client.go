package main

import (
	"context"
	//"fmt"
	log "github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"github.com/therecipe/qt/core"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/connectivity"
	"time"
)

// SprintsClient is GRPC client calling Sprints service
type SprintsClient struct {
	core.QObject

	conn       *grpc.ClientConn
	client     pb.SprintsClient
	tournament *pb.Tournament
	race       *pb.Race

	_ func(msg string)      `signal:"info"`
	_ func(err, msg string) `signal:"error"`
	_ func(msg string)      `signal:"success"`

	_ func(string, uint) string                           `slot:"dialGrpc"`
	_ func(string, uint, rune, uint, []string) string     `slot:"newTournament"`
	_ func([]string, uint) string                         `slot:"newRace"`
	_ func() string                                       `slot:"startRace"`
	_ func() string                                       `slot:"abortRace"`
	_ func(string, string, bool, uint, uint, uint) string `slot:"configureVis"`
}

func SetupSprintsClient(addr string) *SprintsClient {
	client := NewSprintsClient(nil)

	client.ConnectDialGrpc(client.CallDialGrpc)
	client.ConnectNewTournament(client.CallNewTournament)
	client.ConnectNewRace(client.CallNewRace)
	client.ConnectStartRace(client.CallStartRace)
	client.ConnectAbortRace(client.CallAbortRace)
	client.ConnectConfigureVis(client.CallConfigureVis)

	return client
}

func (s *SprintsClient) CallDialGrpc(hostName string, port uint) string {
	var (
		err error
	)

	s.conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(10*time.Second), grpc.WithBlock())
	if err != nil {
		return err.Error()
	}

	s.client = pb.NewSprintsClient(s.conn)

	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//for client.conn.WaitForStateChange(ctx, connState) && connState != connectivity.Ready {
	//connState = client.conn.GetState()
	//}

	//if connState != connectivity.Ready {
	//return nil, fmt.Errorf("couldnt connect to %s in 10 seconds connection state: %s", addr, connState.String())
	//}

	return ""
}

func (s *SprintsClient) CallNewTournament(name string, destValue uint, mode rune, playerCount uint, colors []string) string {
	var err error
	s.tournament, err = s.client.NewTournament(context.Background(), &pb.Tournament{
		Name:        name,
		DestValue:   uint32(destValue),
		Mode:        pb.Tournament_TournamentMode(pb.Tournament_TournamentMode_value[string(mode)]),
		Color:       colors,
		PlayerCount: uint32(playerCount),
	})
	if err != nil {
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) CallNewRace(playerNames []string, destValue uint) string {
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

func (s *SprintsClient) CallStartRace() string {
	_, err := s.client.StartRace(context.Background(), &pb.Empty{})
	if err != nil {
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) CallAbortRace() string {
	_, err := s.client.AbortRace(context.Background(), &pb.Empty{})
	if err != nil {
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) CallConfigureVis(hostName string, visName string, fullscreen bool, resolutionWidth uint, resolutionHeight uint, movingUnit uint) string {
	_, err := s.client.ConfigureVis(context.Background(), &pb.VisConfiguration{
		HostName:         hostName,
		VisName:          visName,
		Fullscreen:       fullscreen,
		ResolutionWidth:  uint32(resolutionWidth),
		ResolutionHeight: uint32(resolutionHeight),
		MovingUnit:       uint32(movingUnit),
	})
	if err != nil {
		return err.Error()
	}
	return ""
}

func (s *SprintsClient) Close() error {
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
