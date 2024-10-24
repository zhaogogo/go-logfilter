package filters

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func NewFilters(filterConfig []any) (*Filters, error) {
	filters := &Filters{
		config: filterConfig,
	}
	for filterIdx, filterC := range filterConfig {
		c := filterC.(map[string]interface{})
		for filterType, filterConfigI := range c {
			log.Info().Msgf("filter[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			filterConfig := filterConfigI.(map[string]any)
			filterplugin, err := GetFilter(filterType, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filter插件不可用 filter[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			}
			filter, err := NewFilter(fmt.Sprintf("%s[%v]", filterType, filterIdx), filterplugin, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filter创建失败 filter[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			}
			filters.filters = append(filters.filters, filter)
		}
	}
	return filters, nil
}

type Filters struct {
	config  []any
	filters []*Filter
}

func (f *Filters) Process(event map[string]interface{}) map[string]interface{} {
	for _, filter := range f.filters {
		if filter.Pass(event) {
			//if filter.prometheusCounter != nil {
			//	filter.prometheusCounter.Inc()
			//}
			event = filter.Process(event)
		}
	}
	return event
}
