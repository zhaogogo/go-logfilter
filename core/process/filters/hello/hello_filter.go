package hello

import (
	"github.com/zhaogogo/go-logfilter/core/topology"
)

type HelloFilter struct {
	name  string
	key   string
	value interface{}
}

func New(config map[string]interface{}) topology.Filter {
	p := &HelloFilter{}
	if v, ok := config["echo"]; ok {
		echo := v.([]interface{})
		if len(echo) < 2 {
			if len(echo) == 0 {
				p.key = ""
				p.value = ""
				return p
			}
			panic("hello 插件 echo 配置列表必须大于二")
		}
		p.key = echo[0].(string)
		p.value = echo[1]
	}

	return p
}

func (p *HelloFilter) Filter(event map[string]any) (map[string]any, error) {
	if p.key != "" && p.value != "" {
		event[p.key] = p.value
	}
	return event, nil
}
