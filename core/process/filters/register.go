package filters

import (
	"fmt"
	"github.com/zhaogogo/go-logfilter/core/process/filters/hello"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"plugin"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func init() {
	Register("hello", hello.New)
}

type BuildFilterFunc func(map[string]interface{}) topology.Processer

var registeredInput map[string]BuildFilterFunc = make(map[string]BuildFilterFunc)

func Register(filterType string, buildFn BuildFilterFunc) {
	if _, ok := registeredInput[filterType]; ok {
		log.Panic().Msgf("filter类型%s已经被注册了", filterType)
	}
	registeredInput[filterType] = buildFn
}

// 获取Filter类型
func GetFilter(filterType string, config map[string]interface{}) (topology.Processer, error) {
	if v, ok := registeredInput[filterType]; ok {
		return v(config), nil
	}
	log.Info().Msgf("filter内置插件[%v]未注册, 尝试加载三方插件", filterType)

	pluginPath := filterType
	filter, err := getFilterFromPlugin(pluginPath, config)
	if err != nil {
		return nil, errors.Wrapf(err, "三方插件%s", pluginPath)
	}
	return filter, nil
}

func getFilterFromPlugin(pluginPath string, config map[string]interface{}) (topology.Processer, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("没有New函数, err=%s", err)
	}
	f, ok := newFunc.(func(map[string]interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("New函数签名错误")
	}
	rst := f(config)
	input, ok := rst.(topology.Processer)
	if !ok {
		return nil, fmt.Errorf("filter未实现Process方法, got: %T", rst)
	}
	return input, nil
}
