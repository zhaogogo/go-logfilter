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
			filterPlugin, err := GetFilter(filterType, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filter插件不可用, filter[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			}
			filterCell, err := NewFilterCell(fmt.Sprintf("%s[%v]", filterType, filterIdx), filterPlugin, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filterCell创建失败input[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			}
			filters.filterCells = append(filters.filterCells, filterCell)
		}
	}
	return filters, nil
}

type Filters struct {
	config      []any
	filterCells []*FilterCell
}

func (f *Filters) Process(event map[string]interface{}) map[string]interface{} {
	for _, filter := range f.filterCells {
		if filter.Pass(event) {
			//if filter.prometheusCounter != nil {
			//	filter.prometheusCounter.Inc()
			//}
			event = filter.Process(event)
		}
	}
	return event
}
