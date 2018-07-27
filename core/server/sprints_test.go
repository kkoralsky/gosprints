package server

import (
	"flag"
	"github.com/golang/mock/gomock"
	"github.com/kkoralsky/gosprints/core/server/mocks"
	"testing"
)

var (
	testSprints *Sprints
	mockDevice  *server_mocks.MockInputDevice
	mockVisMux  *VisMux
)

func TestMain(m *testing.M) {
	flag.Parse()

	// device := NewMockIn

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()
	mockDevice = server_mocks.NewMockInputDevice(ctrl)
	mockVisMux = server_mocks.NewMockVisMuxInterface(ctrl)

	testSprints = SetupSprints(mockDevice, mockVisMux)

	m.Run()
}

func Test_doRace(t *testing.T) {
	testSprints.doRace()
}

func Test_GetResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	stream := server_mocks.NewMockSprints_GetTournamentsServer(ctrl)
	defer ctrl.Finish()

}
