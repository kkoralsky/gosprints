package core

import (
	pb "github.com/kkoralsky/gosprints/proto"
	"log"
	"os"
)

var (
	// ErrorLogger -
	ErrorLogger = log.New(os.Stderr, "E ", log.Lshortfile)
	// InfoLogger -
	InfoLogger = log.New(os.Stdout, "I ", log.Lshortfile)
	// DebugLogger -
	DebugLogger = log.New(os.Stderr, "D ", log.Llongfile)

	DefultTournament = pb.Tournament{
		Name:        "400m distance tournament",
		Color:       []string{"blue", "red", "green", "yellow"},
		Mode:        pb.Tournament_DISTANCE,
		PlayerCount: 2,
		DestValue:   400,
	}
	DummyResults = []*pb.Result{
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 2931.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 23031.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 44231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 1.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 29931.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 2931.3},
		&pb.Result{DestValue: 10, Player: &pb.Player{Gender: pb.Gender(0), Name: "ktos"}, Result: 231.3},
	}
)
