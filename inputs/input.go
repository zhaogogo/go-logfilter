package inputs

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhaogogo/go-logfilter/internal/config"
	"k8s.io/klog/v2"
	"sync"
)

const INPUTS = "inputs"

type Input interface {
	ReadEvent() chan map[string]interface{}
	Shutdown()
}

func NewInputs(appconfig config.AppConfig, config map[string][]interface{}) (inputs *Inputs, err error) {
	inputs = &Inputs{
		appConfig: appconfig,
		config:    config,
	}
	for inputIdx, inputC := range config[INPUTS] {
		var inputPlugin Input
		c := inputC.(map[string]interface{})

		for inputType, inputConfigI := range c {
			klog.Infof("input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI)
			inputConfig := inputConfigI.(map[string]interface{})
			inputPlugin = GetInput(inputType, inputConfig)
			if inputPlugin == nil {
				err = fmt.Errorf("input插件不可用, input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI)
				return
			}
			inputCell, err := NewInputCell(inputPlugin, inputConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "inputCell创建失败input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI)
			}
			inputs.inputCell = append(inputs.inputCell, inputCell)
		}
	}
	return
}

type Inputs struct {
	appConfig config.AppConfig
	config    map[string][]interface{}
	inputCell []*InputCell
}

func (i *Inputs) Start() {
	var wg sync.WaitGroup
	for _, input := range i.inputCell {
		wg.Add(1)
		go func() {
			defer wg.Done()
			input.Start(i.appConfig.Worker)
		}()
	}
	wg.Wait()
}

func (i *Inputs) Stop() {
	for _, input := range i.inputCell {
		input.Shutdown()
	}
}
