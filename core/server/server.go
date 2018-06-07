package server

import (
	"github.com/kkoralsky/gosprints/core"
	"github.com/kkoralsky/gosprints/core/device"
	"github.com/kkoralsky/gosprints/core/visual"
	//"github.com/kkoralsky/gosprints/proto"
)

type CmdServer struct {
	port uint
	//StartTournament(name) error
	//LoadTournament(name) error
	//PrepareRace()
}

func SprintsServer(cfg core.ServerConfig) {
	vis, err := visual.SetupVis(cfg.VisName, cfg.MovingUnit)
	if err != nil {
		panic(err)
	}

	devicePoller, err := device.SetupDevice(cfg.InputDevice, cfg.SamplingRate, cfg.FailstartThreshold)
	if err != nil {
		panic(err)
	}

	go cmdServer.Run()
	go visServer.Run()
	go devicePoller.Run()

	vis.Run()
	devicePoller.Close()
}
