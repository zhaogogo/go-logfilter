package process

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhaogogo/go-logfilter/core/process/filters"
	"github.com/zhaogogo/go-logfilter/core/process/inputs"
	"github.com/zhaogogo/go-logfilter/core/process/outputs"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/internal/config"
	"sync"
)

const (
	INPUTS = "inputs"
	FILTER = "filters"
	OUTPUT = "outputs"
)

type Processer struct {
	ctx       context.Context
	appConfig config.AppConfig
	Inputs    *inputs.Inputs
	processer *Processs
}

func NewProcess(ctx context.Context, appConfig config.AppConfig, conf map[string][]interface{}) (*Processer, error) {
	p := &Processer{
		ctx:       ctx,
		appConfig: appConfig,
		processer: new(Processs),
	}
	if filterConfs, ok := conf[FILTER]; ok {
		if filterConfs == nil || len(filterConfs) == 0 {
		} else {
			filterBox, err := filters.NewFilters(filterConfs)
			if err != nil {
				return nil, err
			}
			p.processer.Add(filterBox)
		}
	}

	if outputConfs, ok := conf[OUTPUT]; ok {
		if outputConfs == nil || len(outputConfs) == 0 {
			return nil, errors.New(fmt.Sprintf("outputs配置插件未配置"))
		}
		outputs, err := outputs.NewOutputs(outputConfs)
		if err != nil {
			return nil, err
		}
		p.processer.Add(outputs)
	} else {
		return nil, errors.New("没有outputs配置")
	}

	if inputConfs, ok := conf[INPUTS]; ok {
		if inputConfs == nil || len(inputConfs) == 0 {
			return nil, errors.New(fmt.Sprintf("inputs配置解析错误, inputs配置: %#v", inputConfs))
		}
		if inputser, err := inputs.NewInputs(inputConfs, p.processer); err != nil {
			return nil, err
		} else {
			p.Inputs = inputser
		}
	} else {
		return nil, errors.New("没有inputs配置")
	}

	return p, nil
}

func (p *Processer) Start() {
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

func (p *Processer) Shutdown() {
	p.Inputs.Stop()
}

type ProcesserNode struct {
	Processer topology.Process
	Next      *ProcesserNode
}

func (p *ProcesserNode) Process(in map[string]interface{}) map[string]interface{} {
	pin := p.Processer.Process(in)
	if p.Next != nil {
		pin = p.Next.Processer.Process(pin)
	}
	return pin
}

type Processs struct {
	process []topology.Process
}

func (p *Processs) Process(in map[string]interface{}) map[string]interface{} {
	for _, pr := range p.process {
		in = pr.Process(in)
	}
	return in
}

func (p *Processs) Add(pr topology.Process) {
	p.process = append(p.process, pr)
}
