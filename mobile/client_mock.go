package main

import (
	gosp "github.com/kkoralsky/gosprints/core"
	pb "github.com/kkoralsky/gosprints/proto"
	"google.golang.org/grpc/connectivity"
)

type mockSprintsClient struct {
	SprintsClient
}

func (m *mockSprintsClient) dialGrpc(string, uint, bool) string {
	m.connState = connectivity.Ready
	m.SetConnState(int(m.connState))
	return ""
}

func (m *mockSprintsClient) newTournament(string, uint, int, uint, []string) string {
	return ""
}

func (m *mockSprintsClient) newRace([]string, uint) string {
	return ""
}

func (m *mockSprintsClient) startRace() string {
	return ""
}

func (m *mockSprintsClient) abortRace() string {
	return ""
}

func (m *mockSprintsClient) configureVis(string, string, bool, uint, uint, uint) string {
	return ""
}

func (m *mockSprintsClient) getResults(string) {
	for _, result := range gosp.DummyResults {
		m.resultModel.AddResult(
			result.Player.Name,
			pb.Gender_name[int32(result.Player.Gender)],
			result.Result,
			uint(result.DestValue),
		)
	}
}

func setupMockSprintsClient(resultModel *ResultModel) *mockSprintsClient {
	client := NewMockSprintsClient(nil)

	client.resultModel = resultModel
	client.connState = connectivity.Shutdown

	client.ConnectDialGrpc(client.dialGrpc)
	client.ConnectNewTournament(client.newTournament)
	client.ConnectNewRace(client.newRace)
	client.ConnectStartRace(client.startRace)
	client.ConnectAbortRace(client.abortRace)
	client.ConnectConfigureVis(client.configureVis)
	client.ConnectGetResults(client.getResults)

	return client
}
