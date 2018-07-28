package visual

import (
	"context"
	"flag"
	"math"
	"testing"
	"time"

	"github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
)

var (
	vis VisInterface
)

func TestMain(m *testing.M) {
	vis_name := flag.String("vis_name", "clock", "either clock or bar")
	flag.Parse()
	vis, _ = SetupVis("gosprints", *vis_name, false, 640, 480, 1, 25*5)
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
		pb.Results{Result: core.DummyResults},
	}

	for _, results := range resultSet {
		vis.ShowResults(context.Background(), &results)
		time.Sleep(2 * time.Second)
	}
}

func Test_PixelVis_Abort(t *testing.T) {
	vis.AbortRace(context.Background(), &pb.AbortMessage{"aborted"})
	time.Sleep(1 * time.Second)
	vis.Clear()
	vis.AbortRace(context.Background(), &pb.AbortMessage{"falsestart"})
	time.Sleep(1 * time.Second)
	vis.Clear()
	vis.AbortRace(context.Background(), &pb.AbortMessage{"fdospf safdso pfs fopsafa fd a"})
	time.Sleep(1 * time.Second)
}

func Test_PixelVis_StartRace(t *testing.T) {
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

func Test_PixelVis_UpdateRace(t *testing.T) {
	vis.NewTournament(context.Background(), &pb.Tournament{
		Color: []string{"blue", "red"}, DestValue: 400, Mode: pb.Tournament_TIME,
		Name: "testing 1 2 3", PlayerCount: 2})
	vis.NewRace(context.Background(), &pb.Race{
		DestValue: 400,
		Players: []*pb.Player{
			&pb.Player{Name: "player1"},
			&pb.Player{Name: "player2"},
		},
	})
	t.Run("test 10sec race", func(tt *testing.T) {
		// maxwait, minwait (miliseconds), maxdelta mindelta, playerCount, maxDistance, maxTime (seconds)
		stream := NewmockRaces_UpdateRaceServer(50, 10, 8, 2, 2, 1000, 10)
		for _, racer := range stream.racers {
			tt.Logf("%d: %d\n", racer.PlayerNum, racer.Distance)
		}
		vis.UpdateRace(stream)
	})
	time.Sleep(time.Second)
	vis.Clear()
	t.Run("test 500 dist race", func(tt *testing.T) {
		stream := NewmockRaces_UpdateRaceServer(50, 10, 8, 2, 2, 500, math.MaxInt32)
		for _, racer := range stream.racers {
			tt.Logf("%d: %d\n", racer.PlayerNum, racer.Distance)
		}
		vis.UpdateRace(stream)
	})
	time.Sleep(time.Second)
	// vis.Clear()
}

func TestFinishRace(t *testing.T) {
	_, err := vis.FinishRace(context.Background(), &pb.Results{
		Result: []*pb.Result{
			{
				Player: &pb.Player{Name: "player1"},
				Result: float32(2442.0 * math.Pow10(9)), // need to scale, because its interpreted as nanosec
			},
			{
				Player: &pb.Player{Name: "player2"},
				Result: float32(3223.4 * math.Pow10(9)),
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

	stream := NewmockRaces_UpdateRaceServer(0, 0, 50, 10, 2, 4, math.MaxInt32)
	b.Run("4 clock moves", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			vis.UpdateRace(stream)
			stream.currentRacer = 0
		}
	})
}
