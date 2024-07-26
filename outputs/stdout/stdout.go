package stdout

import (
	"fmt"
	"github.com/zhaogogo/go-logfilter/outputs"
	"github.com/zhaogogo/go-logfilter/textcodec"
	"k8s.io/klog/v2"
)

func init() {
	outputs.Register("stdout", newStdoutOutput)
}

type StdoutOutput struct {
	config  map[string]interface{}
	encoder textcodec.Encoder
}

func newStdoutOutput(config map[string]interface{}) outputs.Output {
	p := &StdoutOutput{
		config: config,
	}

	if v, ok := config["codec"]; ok {
		p.encoder = textcodec.NewEncoder(v.(string))
	} else {
		p.encoder = textcodec.NewEncoder("json")
	}

	return p

}

func (p *StdoutOutput) Emit(event map[string]interface{}) {
	buf, err := p.encoder.Encode(event)
	if err != nil {
		klog.Errorf("marshal %v error:%s", event, err)
	}
	fmt.Println(string(buf))
}

func (p *StdoutOutput) Shutdown() {}
