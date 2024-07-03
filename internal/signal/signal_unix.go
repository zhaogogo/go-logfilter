package signal

import (
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"os/signal"
	"syscall"
)

func ListenSignal(termFunc func(), reloadFunc func(), cpuProfileFunc func(), memProfileFinc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)

	for sig := range c {
		klog.Infof("capture signal: %v", sig)
		fmt.Println(sig)
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
