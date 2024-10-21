package outputs

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/condition"
	"github.com/zhaogogo/go-logfilter/metrics"
)

func NewOutputCell(outputType string, output Output, cellConfig map[string]interface{}) (*OutputCell, error) {
	o := &OutputCell{
		name:        outputType,
		output:      output,
		config:      cellConfig,
		Conditioner: condition.NewConditioner(cellConfig),
	}
	p, err := metrics.NewPrometheusCounter(cellConfig)
	if err != nil {
		log.Fatal().Err(err)
	}
	o.prometheusCounter = p
	//if add_fields, ok := cellConfig["add_fields"]; ok && add_fields != nil {
	//	o.addFields = make(map[field.FieldSetter]field.ValueRender)
	//	for k, v := range add_fields.(map[string]interface{}) {
	//		fieldSetter := field.NewFieldSetter(k)
	//		if fieldSetter == nil {
	//			klog.Fatalf("fieldSetter构建失败", k)
	//		}
	//		i.addFields[fieldSetter] = field.GetValueRender(v)
	//	}
	//} else {
	//	i.addFields = nil
	//}

	return o, nil
}

type OutputCell struct {
	name              string
	output            Output
	config            map[string]interface{}
	prometheusCounter prometheus.Counter
	exit              func()
	*condition.Conditioner
}

func (o *OutputCell) Emit(event map[string]interface{}) {
	o.output.Emit(event)
}
func (o *OutputCell) Shutdown() {
	//TODO
}
