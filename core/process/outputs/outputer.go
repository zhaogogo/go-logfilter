package outputs

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/condition"
	"github.com/zhaogogo/go-logfilter/pkg/field"
	"github.com/zhaogogo/go-logfilter/pkg/metrics"
)

func NewOutputer(outputType string, output topology.Output, cellConfig map[string]interface{}) (*Outputer, error) {
	o := &Outputer{
		output:      output,
		name:        outputType,
		config:      cellConfig,
		Conditioner: condition.NewConditioner(cellConfig),
		addFields:   field.NewAddFields(cellConfig),
	}
	p, err := metrics.NewPrometheusCounter(cellConfig)
	if err != nil {
		log.Fatal().Err(err)
	}
	o.prometheusCounter = p
	return o, nil
}

type Outputer struct {
	output            topology.Output
	name              string
	config            map[string]interface{}
	prometheusCounter prometheus.Counter
	exit              func()
	*condition.Conditioner

	addFields map[field.FieldSetter]field.ValueRender
}

func (o *Outputer) Emit(event map[string]interface{}) {
	if o.addFields != nil {
		for fs, v := range o.addFields {
			event = fs.SetField(event, v.Render(event), false)
		}
	}
	o.output.Emit(event)
}
