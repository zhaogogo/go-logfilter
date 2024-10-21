package outputs

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func NewOutputs(outputConfig []any) (*Outputs, error) {
	outputs := &Outputs{
		config: outputConfig,
	}
	for outputIdx, outputC := range outputConfig {
		c := outputC.(map[string]interface{})
		for outputType, outputConfigI := range c {
			log.Info().Msgf("output[%d] type: %v config:[%T] %v", outputIdx, outputType, outputConfigI, outputConfigI)
			outputConfig := outputConfigI.(map[string]any)
			outputPlugin, err := GetOutput(outputType, outputConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "output插件不可用, output[%d] type: %v config:[%T] %v", outputIdx, outputType, outputConfigI, outputConfigI)
			}
			outputCell, err := NewOutputCell(fmt.Sprintf("%s[%v]", outputType, outputIdx), outputPlugin, outputConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "outputCell创建失败input[%d] type: %v config:[%T] %v", outputIdx, outputType, outputConfigI, outputConfigI)
			}
			outputs.outputCell = append(outputs.outputCell, outputCell)
		}
	}
	return outputs, nil
}

type Outputs struct {
	config     []any
	outputCell []*OutputCell
}

func (o *Outputs) Process(event map[string]interface{}) map[string]interface{} {
	for _, output := range o.outputCell {
		if output.Pass(event) {
			if output.prometheusCounter != nil {
				output.prometheusCounter.Inc()
			}
			output.Emit(event)
		}
	}
	return nil
}
