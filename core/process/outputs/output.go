package outputs

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/condition"
	"github.com/zhaogogo/go-logfilter/pkg/field"
	"github.com/zhaogogo/go-logfilter/pkg/metrics"
)

func NewOutputer(name string, output topology.Output, cellConfig map[string]interface{}) (*Outputer, error) {
	//var failedtag bool = false
	//failedtagAny, ok := cellConfig["failed_tag"]
	//if ok {
	//	if failedtagBool, ok := failedtagAny.(bool); ok {
	//		failedtag = failedtagBool
	//	}
	//}
	var overwrite bool = false
	overwriteAny, ok := cellConfig["overwrite"]
	if ok {
		if overwriteBool, ok := overwriteAny.(bool); ok {
			overwrite = overwriteBool
		}
	}
	o := &Outputer{
		output:      output,
		name:        name,
		config:      cellConfig,
		Conditioner: condition.NewConditioner(name, cellConfig),
		addFields:   field.NewAddFields(name, cellConfig),
		overwrite:   overwrite,
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

	overwrite bool
	addFields []map[field.FieldSetter]field.ValueRender
}

func (o *Outputer) Emit(event map[string]interface{}) {

	for _, fs := range o.addFields {
		for fieldsetter, valuerender := range fs {
			event = fieldsetter.SetField(event, valuerender.Render(event), o.overwrite)
		}
	}

	o.output.Emit(event)
}
