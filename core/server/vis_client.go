package server

import (
	"context"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"google.golang.org/grpc"
	"strings"
	"time"
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

func (v *VisClient) StartRace(starter *pb.Starter) error {
	for _, cl := range v.clients {
		go cl.StartRace(context.Background(), starter)
	}
	return nil
}

func (v *VisClient) AbortRace(abortMessage *pb.AbortMessage) error {
	for _, cl := range v.clients {
		go cl.AbortRace(context.Background(), abortMessage)
	}
	return nil
}

func (v *VisClient) ConfigureVis(visCfg *pb.VisConfiguration) error {
	for _, cl := range v.clients {
		go cl.ConfigureVis(context.Background(), visCfg)
	}
	return nil
}

func (v *VisClient) connectionStateUpdater() {
	for i, conn := range v.connections {
		go func() {
			for conn != nil {
				conn.WaitForStateChange(context.Background(), conn.GetState())
				if conn != nil {
					core.InfoLogger.Printf("vis connection: %s state: %s",
						v.addresses[i], conn.GetState().String())
				}
			}

			core.InfoLogger.Printf("vis connectino: %s closed", v.addresses[i])
		}()
	}
}

func SetupVisClient(outputs string) (*VisClient, error) {
	var v = VisClient{
		addresses: strings.Split(outputs, ","),
	}
	for _, addr := range v.addresses {
		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(10*time.Second))
		if err != nil {
			core.ErrorLogger.Printf("error while dialing to %s: %s\n", addr, err.Error())
		} else {
			v.connections = append(v.connections, conn)
			v.clients = append(v.clients, pb.NewSprintsClient(conn))
			core.InfoLogger.Printf("dialed to vis: %s with state: %s\n", addr, conn.GetState().String())
		}
	}
	if len(v.connections) > 0 {
		go v.connectionStateUpdater()
		return &v, nil
	}
	return &v, fmt.Errorf("couldnt connect to none of the outputs: %s", outputs)
}
