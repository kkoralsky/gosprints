package main

import (
	"flag"
	"fmt"
	"github.com/kkoralsky/gosprints/core"
	"github.com/kkoralsky/gosprints/core/client"
	"github.com/kkoralsky/gosprints/core/server"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n%s server|client [-help|other options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "server":
			cfg := core.ServerConfig{}
			core.FlagsetParse(cfg.Setup(), os.Args[2:], cfg.Validate)
			server.SprintsServer(cfg)
		case "client":
			cfg := core.ClientConfig{}
			core.FlagsetParse(cfg.Setup(), os.Args[2:], nil)
			client.SprintsClient(cfg)
		default:
			flag.Usage()
		}
	} else {
		flag.Usage()
	}

}
