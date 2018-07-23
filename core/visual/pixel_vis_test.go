package visual

import (
	"context"
	"flag"
	"testing"
	"time"

	pb "github.com/kkoralsky/gosprints/proto"
)

var (
	vis VisInterface
)

func TestMain(m *testing.M) {
	flag.Parse()
	vis, _ = SetupVis("testing 1 2 3", "clock", false, 640, 480, 1)
	go func() {
		time.Sleep(1100 * time.Millisecond) // wait for visualization to setup
		m.Run()
	}()
	vis.Run()
}

func Test_PixelVis_NewRace(t *testing.T) {
	vis.NewTournament(context.Background(), &pb.Tournament{
		Color: []string{"blue", "red"}, DestValue: 400, Mode: 0,
		Name: "testing 1 2 3", PlayerCount: 2})
	vis.NewRace(context.Background(), &pb.Race{
		DestValue: 400,
		Players: []*pb.Player{
			&pb.Player{Name: "player1"},
			&pb.Player{Name: "player2"},
		},
	})
}

func Test_PixelVis_ShowResults(t *testing.T) {
	var resultSet = []pb.Results{
		pb.Results{Result: []*pb.Result{
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
		},
		},
	}

	for _, results := range resultSet {
		vis.ShowResults(context.Background(), &results)
		time.Sleep(2 * time.Second)
	}
}

func Test_PixelVis_Abort(t *testing.T) {
	vis.AbortRace(context.Background(), &pb.AbortMessage{"aborted"})
	time.Sleep(1 * time.Second)
	vis.AbortRace(context.Background(), &pb.AbortMessage{"falsestart"})
	time.Sleep(1 * time.Second)
	vis.AbortRace(context.Background(), &pb.AbortMessage{"fdospf safdso pfs fopsafa fd a"})
}

func Test_ClockVis_StartRace(t *testing.T) {
	vis.NewTournament(context.Background(), &pb.Tournament{
		Color: []string{"blue", "red"}, DestValue: 400, Mode: 0,
		Name: "testing 1 2 3", PlayerCount: 2})
	vis.NewRace(context.Background(), &pb.Race{
		DestValue: 400,
		Players: []*pb.Player{
			&pb.Player{Name: "koral"},
			&pb.Player{Name: "koral2"},
		},
	})
	vis.StartRace(context.Background(), &pb.Starter{3})
}

func Test_ClockVis_UpdateRace(t *testing.T) {
	vis.NewTournament(context.Background(), &pb.Tournament{
		Color: []string{"blue", "red"}, DestValue: 400, Mode: 0,
		Name: "testing 1 2 3", PlayerCount: 2})
	vis.NewRace(context.Background(), &pb.Race{
		DestValue: 400,
		Players: []*pb.Player{
			&pb.Player{Name: "player1"},
			&pb.Player{Name: "player2"},
		},
	})
	stream := NewmockRaces_UpdateRaceServer(5000, 1000, 50, 20, 2, 550)
	for _, racer := range stream.racers {
		t.Logf("%d: %d\n", racer.PlayerNum, racer.Distance)
	}
	vis.UpdateRace(stream)
}

func TestFinishRace(t *testing.T) {
	_, err := vis.FinishRace(context.Background(), &pb.Results{
		Result: []*pb.Result{
			{
				Player: &pb.Player{Name: "player1"},
				// DestValue: uint32(4334),
				Result: float32(2442.0),
			},
			{
				Player: &pb.Player{Name: "player2"},
				// DestValue: uint32(43322),
				Result: float32(3223.4),
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func Benchmark_ClockVis_UpdateRace(b *testing.B) {
	vis.NewTournament(context.Background(), &pb.Tournament{
		Color: []string{"blue", "red"}, DestValue: 400, Mode: 0,
		Name: "testing 1 2 3", PlayerCount: 2})
	vis.NewRace(context.Background(), &pb.Race{
		DestValue: 400,
		Players: []*pb.Player{
			&pb.Player{Name: "player1"},
			&pb.Player{Name: "player2"},
		},
	})

	stream := NewmockRaces_UpdateRaceServer(0, 0, 50, 10, 2, 4)
	b.Run("4 clock moves", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vis.UpdateRace(stream)
			stream.currentRacer = 0
		}
	})
}
