package inputs

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zhaogogo/go-logfilter/field"
	"github.com/zhaogogo/go-logfilter/metrics"
	"k8s.io/klog/v2"
)

func NewInputCell(input Input, cellConfig map[string]interface{}) *InputCell {
	i := &InputCell{
		input:        input,
		config:       cellConfig,
		stop:         false,
		shutdownChan: make(chan struct{}, 1),
	}
	p, err := metrics.NewPrometheusCounter(cellConfig)
	if err != nil {
		klog.Fatal(err)
	}
	i.prometheusCounter = p
	if add_fields, ok := cellConfig["add_fields"]; ok && add_fields != nil {
		i.addFields = make(map[field.FieldSetter]field.ValueRender)
		for k, v := range add_fields.(map[string]interface{}) {
			fieldSetter := field.NewFieldSetter(k)
			if fieldSetter == nil {
				klog.Fatalf("fieldSetter构建失败", k)
			}
			i.addFields[fieldSetter] = field.GetValueRender(v)
		}
	} else {
		i.addFields = nil
	}

	return i
}

type InputCell struct {
	input             Input
	config            map[string]interface{}
	stop              bool
	shutdownChan      chan struct{}
	shutdownWhenNil   bool
	prometheusCounter prometheus.Counter
	exit              func()

	addFields map[field.FieldSetter]field.ValueRender
}

func (i *InputCell) SetShutdownWhenNil(shutdownWhenNil bool) {
	i.shutdownWhenNil = shutdownWhenNil
}

func (i *InputCell) Start() {
	//var firstNode *topology.ProcessorNode = box.buildTopology(workerIdx)

	var (
		event map[string]interface{}
	)

	for !i.stop {
		event = i.input.ReadOneEvent()
		if i.prometheusCounter != nil {
			i.prometheusCounter.Inc()
		}
		if event == nil {
			klog.V(5).Info("received nil message.")
			if i.stop {
				break
			}
			if i.shutdownWhenNil {
				klog.Info("received nil message. shutdown...")
				i.exit()
				break
			} else {
				continue
			}
		}
		for fs, v := range i.addFields {
			event = fs.SetField(event, v.Render(event), "", false)
		}
		//firstNode.Process(event)
		v, _ := json.Marshal(event)
		fmt.Println(string(v))
	}
}
