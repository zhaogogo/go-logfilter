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
		c, ok := filterC.(map[string]interface{})
		if !ok {
			log.Fatal().Msgf("filters asset failed, got: %T, %#v", filterC, filterC)
		}
		for filterType, filterConfigI := range c {

			log.Info().Msgf("filters %v[%d] %v config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
			filterConfig, ok := filterConfigI.(map[string]any)
			if !ok {
				log.Fatal().Msgf("filters %v[%d] config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
			}
			filterplugin, err := GetFilter(filterType, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filters插件不可用 %v[%d] config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
			}
			filter, err := NewFilter(fmt.Sprintf("%s[%v]", filterType, filterIdx), filterplugin, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filters创建失败 %v[%d] config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
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
