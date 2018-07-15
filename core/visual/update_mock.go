package visual

import (
	pb "github.com/kkoralsky/gosprints/proto"
	"io"
	"math/rand"
	"time"
)

type mockRaces_UpdateRaceServer struct {
	pb.Visual_UpdateRaceServer
	racers       []*pb.Racer
	maxWait      int
	minWait      int
	currentRacer int
}

func NewmockRaces_UpdateRaceServer(
	maxWait, minWait,
	maxDistanceDelta, minDistanceDelta,
	playerCount, racersCount int) *mockRaces_UpdateRaceServer {

	rand.Seed(600)
	var (
		m = &mockRaces_UpdateRaceServer{
			minWait:      minWait,
			maxWait:      maxWait,
			currentRacer: 0,
		}
		prevDistance = 0
		currPlayer   = 0
	)

	for i := 0; i < playerCount; i++ {
		m.racers = append(m.racers, &pb.Racer{
			PlayerNum: uint32(i),
			Distance:  0,
		})
	}

	for i := playerCount; i < racersCount; i++ {
		prevDistance = int(m.racers[i-currPlayer-playerCount].Distance)

		m.racers = append(m.racers, &pb.Racer{
			PlayerNum: uint32(currPlayer),
			Distance:  uint32(prevDistance + rand.Intn(maxDistanceDelta-minDistanceDelta) + minDistanceDelta),
		})

		currPlayer++
		if currPlayer == playerCount {
			currPlayer = 0
		}
	}
	return m
}

func (m *mockRaces_UpdateRaceServer) SendAndClose(*pb.Empty) error {
	return nil
}

func (m *mockRaces_UpdateRaceServer) Recv() (*pb.Racer, error) {
	m.currentRacer++
	if m.maxWait != 0 {
		time.Sleep(time.Duration(rand.Intn(m.maxWait-m.minWait)+m.minWait) * time.Nanosecond)
	}
	if m.currentRacer >= len(m.racers) {
		return nil, io.EOF
	}
	return m.racers[m.currentRacer], nil
}
