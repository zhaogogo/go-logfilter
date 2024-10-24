package stdout

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/textcodec"
)

type StdoutOutput struct {
	config  map[string]interface{}
	encoder textcodec.Encoder
}

func New(config map[string]interface{}) topology.Output {
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
		log.Error().Msgf("marshal %v error:%s", event, err)
		return
	}
	fmt.Println(string(buf))
}
