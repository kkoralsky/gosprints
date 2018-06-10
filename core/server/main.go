package server

import (
	"github.com/kkoralsky/gosprints/core"
	"github.com/kkoralsky/gosprints/core/device"
	"github.com/kkoralsky/gosprints/core/visual"
)

func SprintsServer(cfg core.ServerConfig) {
	vis, err := visual.SetupVis(cfg.VisName, cfg.MovingUnit, cfg.Fullscreen)
	if err != nil {
		panic(err)
	}

	devicePoller, err := device.SetupDevice(cfg.InputDevice, cfg.SamplingRate,
		cfg.FailstartThreshold)
	if err != nil {
		panic(err)
	}

	cmdServer, err := SetupCmdServer(cfg.CmdPort, cfg.GrpcDebug)
	if err != nil {
		panic(err)
	}
	visServer, err := SetupVisServer(cfg.VisPort)
	if err != nil {
		panic(err)
	}

	go cmdServer.Run()
	go visServer.Run()
	if err := devicePoller.Start(); err != nil {
		panic(err)
	}

	vis.Run()
	devicePoller.Close()
}
