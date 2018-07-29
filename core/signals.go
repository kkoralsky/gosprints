package core

import (
	"os"
	"os/signal"
	"syscall"
)

func ExitGracefully(handler func()) {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-signals
		handler()

		InfoLogger.Println("exited gracefully")
		os.Exit(0)
	}()
}
