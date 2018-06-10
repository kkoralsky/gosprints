package server

import (
	"context"
	"net"

	pb "github.com/kkoralsky/gosprints/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CmdServer struct {
	port       uint
	grpcServer *grpc.Server
	pb.SprintsServer
	tcpListener *net.TCPListener
}

func (c *CmdServer) NewTournament(ctx context.Context, tournament *pb.Tournament) (*pb.Tournament, error) {
	println(tournament.Name)
	return tournament, nil
}

func (c *CmdServer) NewRace(ctx context.Context, race *pb.Race) (*pb.Empty, error) {
	for _, p := range race.Players {
		println(p.Name)
	}
	return &pb.Empty{}, nil
}

func (c *CmdServer) Run() {
	c.grpcServer.Serve(c.tcpListener)
}

func SetupCmdServer(port uint, debug bool) (*CmdServer, error) {
	var (
		c   = &CmdServer{port: port, grpcServer: grpc.NewServer()}
		err error
	)
	pb.RegisterSprintsServer(c.grpcServer, c)
	if debug {
		reflection.Register(c.grpcServer)
	} else {
		_ = reflection.Register // to prevent "unused import" compiler complaint
	}
	if c.tcpListener, err = net.ListenTCP("tcp", &net.TCPAddr{Port: int(port)}); err != nil {
		return c, err
	}
	return c, nil
}
