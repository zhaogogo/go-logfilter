package filters

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhaogogo/go-logfilter/core/process/filters/convert"
	"github.com/zhaogogo/go-logfilter/core/process/filters/grok"
	"github.com/zhaogogo/go-logfilter/core/process/filters/hello"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"plugin"

	"github.com/rs/zerolog/log"
)

func init() {
	Register("hello", hello.New)
	Register("convert", convert.New)
	Register("grok", grok.New)
	Register("filters", NewFiltersFilters)
}

type BuildFilterFunc func(map[string]interface{}) topology.Filter

var registeredFilter map[string]BuildFilterFunc = make(map[string]BuildFilterFunc)

func Register(filterType string, buildFn BuildFilterFunc) {
	if _, ok := registeredFilter[filterType]; ok {
		log.Panic().Msgf("filter类型%s已经被注册了", filterType)
	}
	registeredFilter[filterType] = buildFn
}

// 获取Filter类型
func GetFilter(filterType string, config map[string]interface{}) (topology.Filter, error) {
	if v, ok := registeredFilter[filterType]; ok {
		return v(config), nil
	}
	log.Info().Msgf("filter内置插件[%v]未注册, 尝试加载三方插件", filterType)

	pluginPath := filterType
	filter, err := getFilterFromPlugin(pluginPath, config)
	if err != nil {
		return nil, err
	}
	return filter, nil
}

func getFilterFromPlugin(pluginPath string, config map[string]interface{}) (topology.Filter, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		return nil, errors.Wrap(err, "三方插件没有New函数")
	}
	f, ok := newFunc.(func(map[string]interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("三方New函数签名错误")
	}
	rst := f(config)
	input, ok := rst.(topology.Filter)
	if !ok {
		return nil, fmt.Errorf("三方插件未实现Process方法, got: %T", rst)
	}
	return input, nil
}
