package server

import (
	"context"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"google.golang.org/grpc"
	"strings"
)

type VisClient struct {
	addresses   []string
	connections []*grpc.ClientConn
	clients     []pb.SprintsClient
	pb.SprintsClient
}

func (v *VisClient) NewTournament(tournament *pb.Tournament) error {
	for _, cl := range v.clients {
		go cl.NewTournament(context.Background(), tournament)
	}
	return nil
}

func (v *VisClient) NewRace(race *pb.Race) error {
	for _, cl := range v.clients {
		go cl.NewRace(context.Background(), race)
	}
	return nil
}

func (v *VisClient) StartRace() error {
	for _, cl := range v.clients {
		go cl.StartRace(context.Background(), &pb.Empty{})
	}
	return nil
}

func (v *VisClient) AbortRace() error {
	for _, cl := range v.clients {
		go cl.AbortRace(context.Background(), &pb.Empty{})
	}
	return nil
}

func (v *VisClient) ConfigureVis(visCfg *pb.VisConfiguration) error {
	for _, cl := range v.clients {
		go cl.ConfigureVis(context.Background(), visCfg)
	}
	return nil
}

func SetupVisClient(outputs string) (*VisClient, error) {
	var v = VisClient{
		addresses: strings.Split(outputs, ","),
	}
	for _, addr := range v.addresses {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			core.ErrorLogger.Printf("error while dialing to %s: %s\n", addr, err.Error())
		} else {
			v.connections = append(v.connections, conn)
			v.clients = append(v.clients, pb.NewSprintsClient(conn))
			core.InfoLogger.Printf("dialed to vis: %s\n", addr)
		}
	}
	if len(v.connections) > 0 {
		return &v, nil
	}
	return &v, fmt.Errorf("couldnt connect to none of the outputs: %s", outputs)
}
