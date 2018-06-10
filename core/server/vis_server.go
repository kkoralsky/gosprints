package server

import (
//pb "github.com/kkoralsky/gosprints/proto"
)

type VisServer struct {
	port uint
}

func (v *VisServer) Run() {
	println("")
}

func SetupVisServer(port uint) (*VisServer, error) {
	return &VisServer{port: port}, nil
}
