package inputs

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
	"plugin"
)

type BuildInputFunc func(map[string]interface{}) Input

var registeredInput map[string]BuildInputFunc = make(map[string]BuildInputFunc)

func Register(inputType string, buildFn BuildInputFunc) {
	if _, ok := registeredInput[inputType]; ok {
		klog.Errorf("input %s已经被注册了, 忽略 %T", inputType, buildFn)
		return
	}
	registeredInput[inputType] = buildFn
}

// 获取INPUT类型
func GetInput(inputType string, config map[string]interface{}) (Input, error) {
	if v, ok := registeredInput[inputType]; ok {
		return v(config), nil
	}
	klog.V(2).Infof("无法加载内置插件, 尝试加载第三方插件", inputType)

	pluginPath := inputType
	output, err := getInputFromPlugin(pluginPath, config)
	if err != nil {
		return nil, errors.Wrap(err, "加载三方插件错误")
	}
	return output, nil
}

func getInputFromPlugin(pluginPath string, config map[string]interface{}) (Input, error) {
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
	input, ok := rst.(Input)
	if !ok {
		return nil, fmt.Errorf("加载三方插件, New函数返回类型错误(%T)", rst)
	}
	return input, nil
}
