package signal

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal(termFunc func(), reloadFunc func(), cpuProfileFunc func(), memProfileFinc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)

	for sig := range c {
		log.Info().Msgf("capture signal: %v", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			termFunc()
		case syscall.SIGHUP:
			reloadFunc()
		case syscall.SIGUSR1:
			cpuProfileFunc()
		case syscall.SIGUSR2:
			memProfileFinc()
		}
	}
}
