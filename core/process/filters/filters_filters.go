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

func NewFiltersFilters(c map[string]any) topology.Filter {
	f := &FiltersFilter{}
	confSlice, ok := c["filter"].([]any)
	if !ok {
		log.Fatal().Msgf("filters plugin config asset failed, got: %T", c)
	}
	failedTag := false
	if failed_tag, ok := c["failed_tag"]; ok {
		failedTag = failed_tag.(bool)
	}
	filters, err := newFiltersFilters(confSlice, failedTag)
	if err != nil {
		panic(err)
	}
	f.filters = filters
	return f
}

func newFiltersFilters(filterConfig []any, failedTag bool) (*Filters, error) {
	filters := &Filters{
		config: filterConfig,
	}
	for filterIdx, filterC := range filterConfig {
		c := filterC.(map[string]interface{})
		for filterType, filterConfigI := range c {
			log.Info().Msgf("filter filters plugin %v[%d] config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
			filterConfig := filterConfigI.(map[string]any)
			filterplugin, err := GetFilter(filterType, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filters filter plugin 插件不可用 %v[%d] config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
			}
			filter, err := NewFilter(fmt.Sprintf("filters filter %s[%v]", filterType, filterIdx), filterplugin, filterConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "filters filter创建失败 %v[%d] config:[%T] %v", filterType, filterIdx, filterConfigI, filterConfigI)
			}
			filter.failed_tag = failedTag
			filters.filters = append(filters.filters, filter)
		}
	}
	return filters, nil
}
