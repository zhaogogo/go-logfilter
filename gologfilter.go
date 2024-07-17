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
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"time"
)

const appName string = "gologfilter"

var (
	version     string
	hostName, _ = os.Hostname()
	pid         = os.Getpid()
)

var appOpts = &struct {
	config     string
	autoReload bool // 配置文件更新自动重启
	pprof      bool
	pprofAddr  string
	cpuprofile string
	memprofile string
	version    bool

	prometheus string

	exitWhenNil bool
}{}

func init() {
	flag.StringVar(&appOpts.config, "config", appOpts.config, "path to configuration file or directory")
	flag.BoolVar(&appOpts.autoReload, "reload", appOpts.autoReload, "if auto reload while config file changed")

	flag.BoolVar(&appOpts.pprof, "pprof", false, "pprof or not")
	flag.StringVar(&appOpts.pprofAddr, "pprof-address", "127.0.0.1:8899", "default: 127.0.0.1:8899")
	flag.StringVar(&appOpts.cpuprofile, "cpuprofile", fmt.Sprintf("cpu.%v", pid), "write cpu profile to `file`")
	flag.StringVar(&appOpts.memprofile, "memprofile", fmt.Sprintf("mem.%v", pid), "write mem profile to `file`")

	flag.BoolVar(&appOpts.version, "version", false, "print version and exit")

	flag.StringVar(&appOpts.prometheus, "prometheus", "", "address to expose prometheus metrics")

	flag.BoolVar(&appOpts.exitWhenNil, "exit-when-nil", false, "triger gohangout to exit when receive a nil event")

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
	if appOpts.version {
		fmt.Printf("%s version %s\n", appName, version)
		return
	}
	klog.Infof("gologfilter version: %s  pid: %v", version, pid)
	defer klog.Flush()
	if appOpts.prometheus != "" {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			klog.Infof("gologfilter prometheus and pprof listen: %s", appOpts.prometheus)
			err := http.ListenAndServe(appOpts.prometheus, nil)
			if err != nil {
				klog.Fatalf("%w", err)
			}
		}()
	}
	if appOpts.pprof {
		go func() {
			http.ListenAndServe(appOpts.pprofAddr, nil)
		}()
	}
	go signal.ListenSignal(exit, reload, CPUProfile, MemProfile)
	conf, err := config.ParseConfig(
		file.NewSource(appOpts.config),
	)
	if err != nil {
		klog.Fatalf("加载配置文件失败", err)
	}
	confy, _ := yaml.Marshal(conf)
	klog.Infof("合并后配置文件为:\n%s", confy)
	inputs, err := inputs.NewInputs(conf)
	if err != nil {
		klog.Fatalf("构建inputs插件失败, err=%v", err)
	}
	go inputs.Start()
	<-ctx.Done()
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
	cpufd, err = os.Create(path.Join(binDir, appOpts.cpuprofile))
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
	memfd, err = os.Create(path.Join(binDir, appOpts.memprofile))
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
