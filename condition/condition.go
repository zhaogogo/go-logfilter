package condition

type Condition interface {
	Pass(event map[string]interface{}) bool
}

type Conditioner struct {
	conditions []Condition
}

func NewConditioner(config map[string]any) *Conditioner {
	f := &Conditioner{}
	// TODO
	return f
}

func (f *Conditioner) Pass(event map[string]interface{}) bool {
	if f.conditions == nil {
		return true
	}

	for _, c := range f.conditions {
		if !c.Pass(event) {
			return false
		}
	}
	return true
}
