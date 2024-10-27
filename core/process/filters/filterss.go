package filters

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
)

type FiltersFilter struct {
	config  []any
	filters *Filters
}

func (f *FiltersFilter) Filter(event map[string]interface{}) (map[string]interface{}, error) {
	event = f.filters.Process(event)
	return event, nil
}

func NewFiltersFilter(c map[string]any) topology.Filter {
	f := &FiltersFilter{}
	confSlice, ok := c["filter"].([]any)
	if !ok {
		log.Fatal().Msgf("filters plugin config asset failed, got: %T", c)
	}
	filters, err := newFiltersFilter(confSlice, c["overwrite"].(bool), c["failed_tag"].(bool))
	if err != nil {
		panic(err)
	}
	f.filters = filters
	return f
}

func newFiltersFilter(filterConfig []any, overwrite bool, failedTag bool) (*Filters, error) {
	filters := &Filters{
		config: filterConfig,
	}
	for filterIdx, filterC := range filterConfig {
		c := filterC.(map[string]interface{})
		for filterType, filterConfigI := range c {
			log.Info().Msgf("filter filters plugin [%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			filterConfig := filterConfigI.(map[string]any)
			filterplugin, err := GetFilter(filterType, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filters filter plugin 插件不可用 filter[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			}
			filter, err := NewFilter(fmt.Sprintf("filters filter %s[%v]", filterType, filterIdx), filterplugin, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filters filter创建失败 filter[%d] type: %v config:[%T] %v", filterIdx, filterType, filterConfigI, filterConfigI)
			}
			filter.failed_tag = failedTag
			filters.filters = append(filters.filters, filter)
		}
	}
	return filters, nil
}
