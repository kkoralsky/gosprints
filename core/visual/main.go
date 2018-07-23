package visual

import (
	"github.com/kkoralsky/gosprints/core"
)

func VisualServer(cfg core.VisualConfig) {
	vis, _ := SetupVis(cfg.HostName, cfg.VisName, cfg.Fullscreen, cfg.ResolutionWidth,
		cfg.ResolutionHeight, cfg.MovingUnit, cfg.DistFactor)

	for vis != nil {
		visServer, err := SetupVisServer(cfg.Port, cfg.GrpcDebug, vis)
		if err != nil {
			panic(err)
		}

		go visServer.Run()
		vis.Run()

		vis = Reconfigure(vis)
		visServer.Stop()
	}
}
