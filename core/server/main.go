package server

import (
	"github.com/kkoralsky/gosprints/core"
	"github.com/kkoralsky/gosprints/core/device"
)

func SprintsServer(cfg core.ServerConfig) {

	devicePoller, err := device.SetupDevice(
		cfg.InputDevice, cfg.SamplingRate, cfg.FailstartThreshold,
	)
	defer devicePoller.Close()

	if err != nil {
		panic(err)
	}

	visMux, err := SetupVisMux(cfg.OutputVisuals)
	if err != nil {
		panic(err)
	}

	sprintsDb, err := SetupSprintsDb(cfg.DbPath)
	defer sprintsDb.Close()
	if err != nil {
		panic(err)
	}
	cmdServer, err := SetupCmdServer(
		cfg.Port, cfg.GrpcDebug,
		SetupSprints(devicePoller, visMux, sprintsDb),
	)
	if err != nil {
		panic(err)
	}

	if err := devicePoller.Start(); err != nil {
		panic(err)
	}

	core.ExitGracefully(func() {
		devicePoller.Close()
		sprintsDb.Close()
		cmdServer.Stop()
	})

	cmdServer.Run()
}
