package inputs

import (
	"fmt"
	"github.com/zhaogogo/go-logfilter/core/process/inputs/stdin"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"plugin"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func init() {
	Register("stdin", stdin.New)
}

type BuildInputFunc func(map[string]interface{}) topology.Input

var registeredInput map[string]BuildInputFunc = make(map[string]BuildInputFunc)

func Register(inputType string, buildFn BuildInputFunc) {
	if _, ok := registeredInput[inputType]; ok {
		log.Panic().Msgf("input类型%s已经被注册了", inputType)
	}
	registeredInput[inputType] = buildFn
}

// 获取INPUT类型
func GetInput(inputType string, config map[string]interface{}) (topology.Input, error) {
	if v, ok := registeredInput[inputType]; ok {
		return v(config), nil
	}
	log.Warn().Msgf("input内置插件[%v]未注册, 尝试加载三方插件", inputType)

	pluginPath := inputType
	output, err := getInputFromPlugin(pluginPath, config)
	if err != nil {
		return nil, errors.Wrapf(err, "三方插件%s", pluginPath)
	}
	return output, nil
}

func getInputFromPlugin(pluginPath string, config map[string]interface{}) (topology.Input, error) {
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
	input, ok := rst.(topology.Input)
	if !ok {
		return nil, fmt.Errorf("未实现Input方法, got: %T", rst)
	}
	return input, nil
}
