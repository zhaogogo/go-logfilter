package process

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhaogogo/go-logfilter/core"
	"github.com/zhaogogo/go-logfilter/filters"
	"github.com/zhaogogo/go-logfilter/inputs"
	"github.com/zhaogogo/go-logfilter/internal/config"
	"github.com/zhaogogo/go-logfilter/outputs"
	"sync"
)

const (
	INPUTS = "inputs"
	FILTER = "filters"
	OUTPUT = "outputs"
)

type Process struct {
	ctx           context.Context
	appConfig     config.AppConfig
	Inputs        *inputs.Inputs
	ProcesserNode *core.ProcesserNode
}

func NewProcess(ctx context.Context, appConfig config.AppConfig, conf map[string][]interface{}) (*Process, error) {
	p := &Process{
		ctx:       ctx,
		appConfig: appConfig,
	}
	if outputConfs, ok := conf[OUTPUT]; ok {
		if outputConfs == nil || len(outputConfs) == 0 {
			return nil, errors.New(fmt.Sprintf("outputs配置解析错误, outputs配置: %#v", outputConfs))
		}
		outputs, err := outputs.NewOutputs(outputConfs)
		if err != nil {
			return nil, err
		}
		p.ProcesserNode = &core.ProcesserNode{
			Processor: outputs,
			Next:      nil,
		}
	} else {
		return nil, errors.New("没有outputs配置")
	}

	if filterConfs, ok := conf[FILTER]; ok {
		if filterConfs == nil || len(filterConfs) == 0 {
		} else {
			filterBox, err := filters.NewFilters(filterConfs)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("创建Filter失败, filter配置: %#v", filterConfs))
			}
			p.ProcesserNode = &core.ProcesserNode{Processor: filterBox, Next: p.ProcesserNode}
		}
	}

	if inputConfs, ok := conf[INPUTS]; ok {
		if inputConfs == nil || len(inputConfs) == 0 {
			return nil, errors.New(fmt.Sprintf("inputs配置解析错误, inputs配置: %#v", inputConfs))
		}
		if inputser, err := inputs.NewInputs(inputConfs, p.ProcesserNode); err != nil {
			return nil, err
		} else {
			p.Inputs = inputser
		}
	} else {
		return nil, errors.New("没有inputs配置")
	}

	return p, nil
}

func (p *Process) Start() {
	wg := sync.WaitGroup{}
	for i := 0; i < p.appConfig.Worker; i++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			p.Inputs.Start(gid)
		}(i)
	}
	wg.Wait()

}

func (p *Process) Shutdown() {
	p.Inputs.Stop()
}
