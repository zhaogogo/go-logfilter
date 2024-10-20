package process

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhaogogo/go-logfilter/inputs"
	"github.com/zhaogogo/go-logfilter/internal/config"
)

const (
	INPUTS = "inputs"
	FILTER = "filters"
	OUTPUT = "outputs"
)

type Processer interface {
	Process(map[string]interface{}) map[string]interface{}
}

type ProcesserNode struct {
	Processor Processer
	Next      *ProcesserNode
}

type Process struct {
	ctx           context.Context
	appConfig     config.AppConfig
	Inputs        *inputs.Inputs
	ProcesserNode *ProcesserNode
}

func NewProcess(ctx context.Context, appConfig config.AppConfig, conf map[string][]interface{}) (*Process, error) {
	p := &Process{
		ctx:       ctx,
		appConfig: appConfig,
	}
	if inputConfs, ok := conf[INPUTS]; ok {
		if inputConfs == nil || len(inputConfs) == 0 {
			return nil, errors.New(fmt.Sprintf("inputs配置解析错误, inputs配置: %#v", inputConfs))
		}
		if inputser, err := inputs.NewInputs(inputConfs); err != nil {
			return nil, err
		} else {
			p.Inputs = inputser
		}
	} else {
		return nil, errors.New("没有inputs配置")
	}
	//if outputConfs, ok := conf[OUTPUT]; ok {
	//	if outputConfs == nil || len(outputConfs) == 0 {
	//		return nil, errors.New(fmt.Sprintf("outputs配置解析错误, outputs配置: %#v", outputConfs))
	//	}
	//	outputs, err := outputs.NewOutputs(outputConfs)
	//	if err != nil {
	//		return nil, err
	//	}
	//	p.ProcesserNode = &ProcesserNode{
	//		Processor: outputs,
	//		Next:      nil,
	//	}
	//} else {
	//	return nil, errors.New("没有outputs配置")
	//}
	return p, nil
}

func (p *Process) Start() {
	p.Inputs.Start()
}

func (p *Process) Shutdown() {
	p.Inputs.Stop()
}
