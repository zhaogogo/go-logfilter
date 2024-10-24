package filters

import (
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/condition"
	"github.com/zhaogogo/go-logfilter/pkg/field"
)

func NewFilter(filterType string, filter topology.Processer, cellConfig map[string]interface{}) (*Filter, error) {
	f := &Filter{
		name:        filterType,
		filter:      filter,
		config:      cellConfig,
		Conditioner: condition.NewConditioner(cellConfig),
		addFields:   field.NewAddFields(cellConfig),
	}
	return f, nil
}

type Filter struct {
	name   string
	filter topology.Processer
	config map[string]interface{}
	//prometheusCounter prometheus.Counter
	exit func()

	addFields map[field.FieldSetter]field.ValueRender
	*condition.Conditioner
}

func (f *Filter) Process(event map[string]interface{}) map[string]interface{} {
	if f.Conditioner.Pass(event) {
		event = f.filter.Process(event)
		for fs, v := range f.addFields {
			event = fs.SetField(event, v.Render(event), false)
		}
	}

	return event
}
