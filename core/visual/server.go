package visual

import (
	pb "github.com/kkoralsky/gosprints/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type VisServer struct {
	port        uint
	grpcServer  *grpc.Server
	tcpListener *net.TCPListener
	vis         VisInterface
}

func (v *VisServer) Run() {
	v.grpcServer.Serve(v.tcpListener)
}

func (v *VisServer) Stop() {
	v.grpcServer.Stop()
	v.tcpListener.Close()
}

func SetupVisServer(port uint, debug bool, vis VisInterface) (*VisServer, error) {
	var (
		v = &VisServer{
			port:       port,
			grpcServer: grpc.NewServer(),
			vis:        vis,
		}
		err error
	)
	pb.RegisterVisualServer(v.grpcServer, v.vis)

	if debug {
		reflection.Register(v.grpcServer)
	} else {
		_ = reflection.Register // to prevent "unused import" compiler complaint
	}

	if v.tcpListener, err = net.ListenTCP("tcp", &net.TCPAddr{Port: int(port)}); err != nil {
		return v, err
	}
	return v, nil
}
