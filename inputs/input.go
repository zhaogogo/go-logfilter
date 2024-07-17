package inputs

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

const INPUTS = "inputs"

type Input interface {
	ReadOneEvent() map[string]interface{}
	Shutdown()
}

func NewInputs(config map[string][]interface{}) (inputs *Inputs, err error) {
	inputs = &Inputs{
		config: config,
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
			inputCell := NewInputCell(inputPlugin, inputConfig)
			if inputCell == nil {
				err = errors.New(fmt.Sprintf("inputCell创建失败input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI))
				return
			}
			inputs.inputCell = append(inputs.inputCell, inputCell)
		}
	}
	return
}

type Inputs struct {
	config    map[string][]interface{}
	inputCell []*InputCell
}

func (i *Inputs) Start() {
	for _, input := range i.inputCell {
		input.Start()
	}
}
