package inputs

import (
	"fmt"
	"github.com/zhaogogo/go-logfilter/core/topology"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func NewInputs(inputsConf []any, process topology.Process) (inputs *Inputs, err error) {
	inputs = &Inputs{
		config: inputsConf,
	}

	for inputIdx, inputC := range inputsConf {
		var input topology.Input
		c := inputC.(map[string]interface{})

		for inputType, inputConfigI := range c {
			log.Info().Msgf("input-%v[%d] config:[%T] %v", inputType, inputIdx, inputConfigI, inputConfigI)
			inputConfig := inputConfigI.(map[string]interface{})
			input, err = GetInput(inputType, inputConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "input插件不可用, input-%v[%d] type: %v config:[%T] %v", inputType, inputIdx, inputConfigI, inputConfigI)
			}
			inputer, err := NewInputer(fmt.Sprintf("input-%s[%v]", inputType, inputIdx), input, inputConfig, process)
			if err != nil {
				return nil, errors.Wrapf(err, "input创建失败 input-%v[%d] config:[%T] %v", inputType, inputIdx, inputConfigI, inputConfigI)
			}
			inputs.inputs = append(inputs.inputs, inputer)
		}
	}
	return
}

type Inputs struct {
	config []any
	inputs []*Inputer
}

func (i *Inputs) Start(gid int) {
	// 所有inputs开始读取事件
	for _, input := range i.inputs {
		input.Start(gid)
	}
}

func (i *Inputs) Stop() {
	for _, input := range i.inputs {
		input.Shutdown()
	}
}
