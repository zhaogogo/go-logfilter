package filters

import (
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/condition"
	"github.com/zhaogogo/go-logfilter/pkg/field"
)

func NewFilter(filterType string, filter topology.Filter, cellConfig map[string]interface{}) (*Filter, error) {
	var failedtag bool = false
	failedtagAny, ok := cellConfig["failed_tag"]
	if ok {
		if failedtagBool, ok := failedtagAny.(bool); ok {
			failedtag = failedtagBool
		}

	}
	var overwrite bool = false
	overwriteAny, ok := cellConfig["overwrite"]
	if ok {
		if overwriteBool, ok := overwriteAny.(bool); ok {
			failedtag = overwriteBool
		}
	}
	f := &Filter{
		name:         filterType,
		filter:       filter,
		config:       cellConfig,
		failed_tag:   failedtag,
		addFields:    field.NewAddFields(cellConfig, failedtag),
		overwrite:    overwrite,
		deleteFields: field.NewFieldDeleter(cellConfig, filterType),
		Conditioner:  condition.NewConditioner(cellConfig),
	}
	return f, nil
}

type Filter struct {
	name   string
	filter topology.Filter
	config map[string]interface{}
	//prometheusCounter prometheus.Counter
	failed_tag   bool
	overwrite    bool
	addFields    []map[field.FieldSetter]field.ValueRender
	deleteFields []field.FieldDelete
	*condition.Conditioner
}

func (f *Filter) Process(event map[string]interface{}) map[string]interface{} {
	var err error
	if f.Conditioner.Pass(event) {
		event, err = f.filter.Filter(event)
		if err != nil && f.failed_tag {
			event = field.SetFailedTags(event, field.FailedTagKey, err.Error())
		}
		if err == nil {
			for _, fs := range f.addFields {
				for fieldsetter, valuerender := range fs {
					event = fieldsetter.SetField(event, valuerender.Render(event), f.overwrite)
				}
			}

			for _, fdel := range f.deleteFields {
				fdel.Delete(event)
			}
		}
	}

	return event
}
