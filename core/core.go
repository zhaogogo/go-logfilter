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
