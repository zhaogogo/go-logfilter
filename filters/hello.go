package filters

import (
	"fmt"
	"github.com/zhaogogo/go-logfilter/core"
)

type HelloFilter struct {
	echo string
}

func init() {
	Register("hello", newHelloFilter)
}

func newHelloFilter(config map[string]interface{}) core.Processer {
	echo := "hello"
	if v, ok := config["echo"]; ok {
		echo = v.(string)
	}
	p := &HelloFilter{echo: echo}

	return p
}

func (p *HelloFilter) Process(event map[string]any) map[string]any {
	if p.echo != "" {
		fmt.Println(p.echo)
	} else {
		fmt.Println("hello filter")
	}
	return event
}
