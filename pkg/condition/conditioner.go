package condition

import "github.com/rs/zerolog/log"

type Condition interface {
	Pass(event map[string]interface{}) bool
}

type Conditioner struct {
	conditions []Condition
}

func NewConditioner(name string, config map[string]any) *Conditioner {
	f := &Conditioner{}
	if v, ok := config["if"]; ok {
		cc, ok := v.([]any)
		if !ok {
			log.Panic().Msgf("%s condition if config assert []any incorrect， got %T", name, v)
		}
		f.conditions = make([]Condition, len(cc))
		for i, c := range cc {
			f.conditions[i] = NewCondition(name, c.(string))
		}
	} else {
		f.conditions = nil
	}
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
