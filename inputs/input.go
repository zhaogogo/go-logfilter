package inputs

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Input interface {
	ReadEvent() chan map[string]interface{}
	Shutdown()
}

func NewInputs(inputsConf []any) (inputs *Inputs, err error) {
	inputs = &Inputs{
		config: inputsConf,
	}
	for inputIdx, inputC := range inputsConf {
		var inputPlugin Input
		c := inputC.(map[string]interface{})

		for inputType, inputConfigI := range c {
			log.Info().Msgf("input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI)
			inputConfig := inputConfigI.(map[string]interface{})
			inputPlugin, err = GetInput(inputType, inputConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "input插件不可用, input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI)
			}
			inputCell, err := NewInputCell(fmt.Sprintf("%s[%v]", inputType, inputIdx), inputPlugin, inputConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "inputCell创建失败input[%d] type: %v config:[%T] %v", inputIdx, inputType, inputConfigI, inputConfigI)
			}
			inputs.inputCell = append(inputs.inputCell, inputCell)
		}
	}
	return
}

type Inputs struct {
	config    []any
	inputCell []*InputCell
}

func (i *Inputs) Start() {
	// 所有inputs开始读取事件
	for _, input := range i.inputCell {
		input.Start()
	}
}

func (i *Inputs) Stop() {
	for _, input := range i.inputCell {
		input.Shutdown()
	}
}
