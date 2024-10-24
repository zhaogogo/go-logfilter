package core

type Processer interface {
	Process(map[string]interface{}) map[string]interface{}
}

type ProcesserNode struct {
	Processor Processer
	Next      *ProcesserNode
}

func (p *ProcesserNode) Process(in map[string]interface{}) map[string]interface{} {
	pin := p.Processor.Process(in)
	if p.Next != nil {
		pin = p.Next.Processor.Process(pin)
	}
	return pin
}

type Process = []Processer

func (p *Process) Process(in map[string]interface{}) map[string]interface{} {
	for _, pr := range *p {
		in = pr.Process(in)
	}
	return in
}
