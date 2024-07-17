package inputs

import (
	"fmt"
	"k8s.io/klog/v2"
	"plugin"
)

type BuildInputFunc func(map[string]interface{}) Input

var registeredInput map[string]BuildInputFunc = make(map[string]BuildInputFunc)

func Register(inputType string, buildFn BuildInputFunc) {
	if _, ok := registeredInput[inputType]; ok {
		klog.Errorf("%s已经被注册了, 忽略 %T", inputType, buildFn)
		return
	}
	registeredInput[inputType] = buildFn
}

// 获取INPUT类型
func GetInput(inputType string, config map[string]interface{}) Input {
	if v, ok := registeredInput[inputType]; ok {
		return v(config)
	}
	klog.Infof("无法加载input[%s] %v插件, 尝试加载第三方插件", inputType, config)

	pluginPath := inputType
	output, err := getInputFromPlugin(pluginPath, config)
	if err != nil {
		klog.Errorf("加载三方插件错误, err=%v", err)
		return nil
	}
	return output
}

func getInputFromPlugin(pluginPath string, config map[string]interface{}) (Input, error) {
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("加载插件%s %v, err=%v", pluginPath, config, err)
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		return nil, fmt.Errorf("加载插件%s %v, 没有`New`函数, err=%s", pluginPath, config, err)
	}
	f, ok := newFunc.(func(map[string]interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("加载插件%s %v, `New`函数签名错误", pluginPath, config)
	}
	rst := f(config)
	input, ok := rst.(Input)
	if !ok {
		return nil, fmt.Errorf("加载插件%s, `New`函数返回类型错误(%T)", pluginPath, rst)
	}
	return input, nil
}
