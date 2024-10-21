package textcodec

import (
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/textcodec/json"
	"github.com/zhaogogo/go-logfilter/textcodec/plain"
	"plugin"
)

type Decoder interface {
	Decode([]byte) map[string]interface{}
}

func NewDecoder(t string, timestampKey string) Decoder {
	switch t {
	case "plain":
		return &plain.PlainDecoder{TimestampKey: timestampKey}
	case "json":
		return &json.Jsoner{TimestampKey: timestampKey, UseNumber: true}
	case "json:not_usenumber":
		return &json.Jsoner{TimestampKey: timestampKey, UseNumber: false}
	default:
		p, err := plugin.Open(t)
		if err != nil {
			log.Fatal().Msgf("could not open %s: %s", t, err)
		}
		newFunc, err := p.Lookup("New")
		if err != nil {
			log.Fatal().Msgf("could not find New function in %s: %s", t, err)
		}
		return newFunc.(func() interface{})().(Decoder)
	}
}
