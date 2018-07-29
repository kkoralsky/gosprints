package server

import (
	"net"

	pb "github.com/kkoralsky/gosprints/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CmdServer struct {
	port        uint
	grpcServer  *grpc.Server
	Sprints     *Sprints
	tcpListener *net.TCPListener
}

func (c *CmdServer) Run() {
	c.grpcServer.Serve(c.tcpListener)
}

func (c *CmdServer) Stop() {
	c.grpcServer.Stop()
	c.tcpListener.Close()
}

func SetupCmdServer(port uint, debug bool, sprints *Sprints) (*CmdServer, error) {
	var (
		c = &CmdServer{
			port:       port,
			grpcServer: grpc.NewServer(),
			Sprints:    sprints,
		}
		err error
	)
	pb.RegisterSprintsServer(c.grpcServer, c.Sprints)
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
