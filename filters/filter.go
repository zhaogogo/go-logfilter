package filters

import (
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/condition"
	"github.com/zhaogogo/go-logfilter/core"
	"github.com/zhaogogo/go-logfilter/field"
)

func NewFilterCell(filterType string, filter core.Processer, cellConfig map[string]interface{}) (*FilterCell, error) {
	f := &FilterCell{
		name:        filterType,
		filter:      filter,
		config:      cellConfig,
		Conditioner: condition.NewConditioner(cellConfig),
	}
	//p, err := metrics.NewPrometheusCounter(cellConfig)
	//if err != nil {
	//	log.Fatal().Err(err)
	//}
	//f.prometheusCounter = p
	if add_fields, ok := cellConfig["add_fields"]; ok && add_fields != nil {
		f.addFields = make(map[field.FieldSetter]field.ValueRender)
		for k, v := range add_fields.(map[string]interface{}) {
			fieldSetter := field.NewFieldSetter(k)
			if fieldSetter == nil {
				log.Fatal().Msgf("filter fieldSetter构建失败", k)
			}
			f.addFields[fieldSetter] = field.GetValueRender(v)
		}
	} else {
		f.addFields = nil
	}

	return f, nil
}

type FilterCell struct {
	name   string
	filter core.Processer
	config map[string]interface{}
	//prometheusCounter prometheus.Counter
	exit func()

	addFields map[field.FieldSetter]field.ValueRender
	*condition.Conditioner
}

func (f *FilterCell) Process(event map[string]interface{}) map[string]interface{} {
	if f.Conditioner.Pass(event) {
		event = f.filter.Process(event)
		for fs, v := range f.addFields {
			event = fs.SetField(event, v.Render(event), "", false)
		}
	}

	return event
}
