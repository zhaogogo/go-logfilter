package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhaogogo/go-logfilter/inputs"
	"github.com/zhaogogo/go-logfilter/internal/config"
	"github.com/zhaogogo/go-logfilter/internal/config/file"
	"github.com/zhaogogo/go-logfilter/internal/signal"
	"k8s.io/klog/v2"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"time"
)

const appName string = "logfilter"

var (
	version     string
	hostName, _ = os.Hostname()
	pid         = os.Getpid()
)

var appOpts = config.AppConfig{
	Config:      "config",
	AutoReload:  false,
	Pprof:       false,
	PprofAddr:   "",
	Cpuprofile:  fmt.Sprintf("cpu.%v", pid),
	Memprofile:  fmt.Sprintf("mem.%v", pid),
	Version:     false,
	Prometheus:  "",
	ExitWhenNil: false,
	Worker:      1,
}

func init() {
	flag.StringVar(&appOpts.Config, "config", appOpts.Config, "path to configuration file or directory")
	flag.BoolVar(&appOpts.AutoReload, "reload", appOpts.AutoReload, "if auto reload while config file changed")

	flag.BoolVar(&appOpts.Pprof, "pprof", false, "pprof or not")
	flag.StringVar(&appOpts.PprofAddr, "pprof-address", "127.0.0.1:8899", "default: 127.0.0.1:8899")
	flag.StringVar(&appOpts.Cpuprofile, "cpuprofile", fmt.Sprintf("cpu.%v", pid), "write cpu profile to `file`")
	flag.StringVar(&appOpts.Memprofile, "memprofile", fmt.Sprintf("mem.%v", pid), "write mem profile to `file`")

	flag.BoolVar(&appOpts.Version, "version", false, "print version and exit")
	flag.StringVar(&appOpts.Prometheus, "prometheus", "", "address to expose prometheus metrics")

	flag.BoolVar(&appOpts.ExitWhenNil, "exit-when-nil", false, "triger gohangout to exit when receive a nil event")
	flag.IntVar(&appOpts.Worker, "worker", 1, "worker thread count")
	klog.InitFlags(nil)
	flag.Parse()
}

var (
	ctx    context.Context
	cancel context.CancelFunc

	cpufd *os.File
	memfd *os.File
)

// TODO
func reload() {
}

func main() {
	flag.Parse()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	if appOpts.Version {
		fmt.Printf("%s version %s\n", appName, version)
		return
	}
	klog.Infof("gologfilter version: %s  pid: %v", version, pid)
	defer klog.Flush()
	if appOpts.Prometheus != "" {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			klog.Infof("gologfilter prometheus and pprof listen: %s", appOpts.Prometheus)
			err := http.ListenAndServe(appOpts.Prometheus, nil)
			if err != nil {
				klog.Fatalf("%w", err)
			}
		}()
	}
	if appOpts.Pprof {
		go func() {
			http.ListenAndServe(appOpts.PprofAddr, nil)
		}()
	}
	go signal.ListenSignal(exit, reload, CPUProfile, MemProfile)
	conf, err := config.ParseConfig(
		file.NewSource(appOpts.Config),
	)
	if err != nil {
		klog.Fatalf("加载配置文件失败", err)
	}

	inputs, err := inputs.NewInputs(appOpts, conf)
	if err != nil {
		klog.Fatalf("构建inputs插件失败, err=%v", err)
	}
	go func() {
		<-ctx.Done()
		fmt.Println("inputs.Stop()")
		inputs.Stop()
	}()
	inputs.Start()
}

func exit() {
	cancel()
}

func CPUProfile() {
	if cpufd == nil {
		cpuProfileStart()
		go func() {
			time.Sleep(time.Second * 30)
			cpuProfileStop()
		}()
	} else {
		cpuProfileStop()
	}

}

func cpuProfileStart() {
	klog.Infof("开始收集CPU信息cpu.%v profile文件", pid)
	binPath, err := os.Executable()
	if err != nil {
		klog.Infof("创建CPU profile, 获取进程执行路径失败: %s", err)
		return
	}
	binDir := path.Dir(binPath)
	cpufd, err = os.Create(path.Join(binDir, appOpts.Cpuprofile))
	if err != nil {
		klog.Infof("could not create CPU profile: %s", err)
		return
	}
	if err := pprof.StartCPUProfile(cpufd); err != nil {
		klog.Infof("could not start CPU profile: %s", err)
		return
	}
}

func cpuProfileStop() {
	klog.Infof("结束收集CPU信息mem.%v profile文件", pid)
	if cpufd == nil {
		klog.Warning("could not close CPU profile file FD: FD is nill, 可能已经手动停止")
		return
	}
	pprof.StopCPUProfile()
	if err := cpufd.Close(); err != nil {
		klog.Fatalf("could not close CPU profile file FD: %s", err)
	}
	cpufd = nil
}

func MemProfile() {
	if memfd == nil {
		memProfileStart()
		go func() {
			time.Sleep(time.Second * 30)
			memProfileStop()
		}()
	} else {
		memProfileStop()
	}
}

func memProfileStart() {
	klog.Infof("开始收集内存信息mem.%v profile文件", pid)
	binPath, err := os.Executable()
	if err != nil {
		klog.Infof("创建MEM profile, 获取进程执行路径失败: %s", err)
		return
	}
	binDir := path.Dir(binPath)
	memfd, err = os.Create(path.Join(binDir, appOpts.Memprofile))
	if err != nil {
		klog.Infof("could not create memory profile: %s", err)
		return
	}

	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(memfd); err != nil {
		klog.Infof("could not write memory profile: %s", err)
		if err := memfd.Close(); err != nil {
			klog.Fatalf("could not close memory profile file FD: %s", err)
		}
		return
	}
}

func memProfileStop() {
	klog.Infof("结束收集内存信息mem.%v profile文件", pid)
	if memfd == nil {
		klog.Warning("could not close memory profile file FD: FD is nill, 可能已经手动停止")
		return
	}
	if err := memfd.Close(); err != nil {
		klog.Fatalf("could not close memory profile file FD: %s", err)
	}
	memfd = nil
}
