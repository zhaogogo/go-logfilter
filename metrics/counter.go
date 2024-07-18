package metrics

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"k8s.io/klog/v2"
	"sync"
)

var lock = sync.Mutex{}
var counterManager map[string]prometheus.Counter = make(map[string]prometheus.Counter)

func hashValue(opts prometheus.CounterOpts) string {
	opts.Help = ""
	b, _ := json.Marshal(opts)
	return string(b)
}

func NewPrometheusCounter(config map[string]interface{}) (prometheus.Counter, error) {
	if promConf, ok := config["prometheus_counter"]; ok {
		var opts prometheus.CounterOpts = prometheus.CounterOpts{}
		err := mapstructure.Decode(promConf, &opts)
		if err != nil {
			return nil, errors.Wrap(err, "prometheusOpts配置解析失败")
		}

		key := hashValue(opts)
		klog.Infof("创建prometheus指标, %s", key)
		if _, ok := counterManager[key]; ok {
			return nil, errors.New(fmt.Sprintf("prometheus配置重复 prometheus_conf=%v", promConf))
		}
		c := promauto.NewCounter(opts)
		lock.Lock()
		counterManager[key] = c
		lock.Unlock()
		return c, nil
	}
	return nil, nil
}
