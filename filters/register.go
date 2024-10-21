package filters

import (
	"fmt"
	"github.com/zhaogogo/go-logfilter/core"
	"plugin"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type BuildFilterFunc func(map[string]interface{}) core.Processer

var registeredInput map[string]BuildFilterFunc = make(map[string]BuildFilterFunc)

func Register(filterType string, buildFn BuildFilterFunc) {
	if _, ok := registeredInput[filterType]; ok {
		log.Error().Msgf("filter %s已经被注册了, 忽略 %T", filterType, buildFn)
		return
	}
	registeredInput[filterType] = buildFn
}

// 获取Filter类型
func GetFilter(filterType string, config map[string]interface{}) (core.Processer, error) {
	if v, ok := registeredInput[filterType]; ok {
		return v(config), nil
	}
	log.Info().Msgf("无法加载内置插件, 尝试加载第三方插件, %v", filterType)

	pluginPath := filterType
	filter, err := getFilterFromPlugin(pluginPath, config)
	if err != nil {
		return nil, errors.Wrap(err, "加载三方插件错误")
	}
	return filter, nil
}

func getFilterFromPlugin(pluginPath string, config map[string]interface{}) (core.Processer, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("加载三方插件, 没有New函数, err=%s", err)
	}
	f, ok := newFunc.(func(map[string]interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("加载三方插件, New函数签名错误")
	}
	rst := f(config)
	input, ok := rst.(core.Processer)
	if !ok {
		return nil, fmt.Errorf("加载三方插件, New函数返回类型错误(%T)", rst)
	}
	return input, nil
}
