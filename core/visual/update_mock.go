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
	finishTime   time.Time
}

func NewmockRaces_UpdateRaceServer(
	maxWait, minWait,
	maxDistanceDelta, minDistanceDelta,
	playerCount, maxDistance, maxTime int) *mockRaces_UpdateRaceServer {

	rand.Seed(600)
	var (
		m = &mockRaces_UpdateRaceServer{
			minWait:      minWait,
			maxWait:      maxWait,
			currentRacer: 0,
			finishTime:   time.Now().Add(time.Duration(maxTime) * time.Second),
		}
		prevDistance   uint32 = 0
		currDistance   uint32 = 0
		currPlayer     uint32 = 0
		playerFinished uint32 = 0
	)

	for i := 0; i < playerCount; i++ {
		m.racers = append(m.racers, &pb.Racer{
			PlayerNum: uint32(i),
			Distance:  0,
		})
	}

	for i := playerCount; playerFinished < uint32(playerCount); i++ {
		prevDistance = m.racers[i-int(currPlayer)-playerCount].Distance
		currDistance = prevDistance + uint32(rand.Intn(maxDistanceDelta-minDistanceDelta)+minDistanceDelta)

		m.racers = append(m.racers, &pb.Racer{
			PlayerNum: currPlayer,
			Distance:  currDistance,
		})
		if currDistance > uint32(maxDistance) {
			playerFinished++
		}

		currPlayer++
		if currPlayer == uint32(playerCount) {
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
		time.Sleep(time.Duration(rand.Intn(m.maxWait-m.minWait)+m.minWait) * time.Millisecond)
	}
	if m.currentRacer >= len(m.racers) || time.Now().After(m.finishTime) {
		return nil, io.EOF
	}
	return m.racers[m.currentRacer], nil
}
