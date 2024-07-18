package config

type AppConfig struct {
	Config     string
	AutoReload bool // 配置文件更新自动重启
	Pprof      bool
	PprofAddr  string
	Cpuprofile string
	Memprofile string
	Version    bool

	Prometheus string

	ExitWhenNil bool
	Worker      int
}
