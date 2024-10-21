package outputs

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"plugin"
)

type Output interface {
	Emit(map[string]interface{})
	Shutdown()
}

type BuildOutputFunc func(map[string]interface{}) Output

var registeredOutput map[string]BuildOutputFunc = make(map[string]BuildOutputFunc)

func Register(outputType string, buildFn BuildOutputFunc) {
	if _, ok := registeredOutput[outputType]; ok {
		log.Error().Msgf("output %s已经被注册了, 忽略 %T", outputType, buildFn)
		return
	}
	registeredOutput[outputType] = buildFn
}

// 获取OUTPUT类型
func GetOutput(outputType string, config map[string]interface{}) (Output, error) {
	if v, ok := registeredOutput[outputType]; ok {
		return v(config), nil
	}
	log.Info().Msgf("无法加载output[%s]插件, 尝试加载第三方插件", outputType)

	pluginPath := outputType
	output, err := getOutputFromPlugin(pluginPath, config)
	if err != nil {
		return nil, errors.Wrap(err, "加载三方插件错误")
	}
	return output, nil
}

func getOutputFromPlugin(pluginPath string, config map[string]interface{}) (Output, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, err
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("加载插件%s, 没有New函数, err=%s", pluginPath, err)
	}
	f, ok := newFunc.(func(map[string]interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("加载插件%s, New函数签名错误", pluginPath)
	}
	rst := f(config)
	input, ok := rst.(Output)
	if !ok {
		return nil, fmt.Errorf("加载插件%s, New函数返回类型错误(%T)", pluginPath, rst)
	}
	return input, nil
}
