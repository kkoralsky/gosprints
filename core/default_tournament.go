package core

import (
	pb "github.com/kkoralsky/gosprints/proto"
)

var DefultTournament = pb.Tournament{
	Name:        "400m distance tournament",
	Color:       []string{"blue", "red", "green", "yellow"},
	Mode:        pb.Tournament_DISTANCE,
	PlayerCount: 2,
	DestValue:   400,
}
