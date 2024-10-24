package outputs

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/process/outputs/stdout"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"plugin"
)

func init() {
	Register("stdout", stdout.New)
}

type BuildOutputFunc func(map[string]interface{}) topology.Output

var registeredOutput map[string]BuildOutputFunc = make(map[string]BuildOutputFunc)

func Register(outputType string, buildFn BuildOutputFunc) {
	if _, ok := registeredOutput[outputType]; ok {
		log.Panic().Msgf("output类型%s已经被注册了", outputType)
	}
	registeredOutput[outputType] = buildFn
}

// 获取OUTPUT类型
func GetOutput(outputType string, config map[string]interface{}) (topology.Output, error) {
	if v, ok := registeredOutput[outputType]; ok {
		return v(config), nil
	}
	log.Warn().Msgf("output内置插件[%v]未注册, 尝试加载三方插件", outputType)

	pluginPath := outputType
	output, err := getOutputFromPlugin(pluginPath, config)
	if err != nil {
		return nil, errors.Wrapf(err, "三方插件%s", pluginPath)
	}
	return output, nil
}

func getOutputFromPlugin(pluginPath string, config map[string]interface{}) (topology.Output, error) {
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
	input, ok := rst.(topology.Output)
	if !ok {
		return nil, fmt.Errorf("未实现Output方法, got: %T", rst)
	}
	return input, nil
}
