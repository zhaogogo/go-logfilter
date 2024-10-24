package textcodec

import (
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/pkg/textcodec/json"
	"plugin"
)

type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

func NewEncoder(t string) Encoder {
	switch t {
	case "json":
		return &json.Jsoner{UseNumber: true}
		//case "simplejson":
		//	return &simplejson.SimpleJsonDecoder{}
	}

	// try plugin
	p, err := plugin.Open(t)
	if err != nil {
		log.Fatal().Msgf("could not open %s: %s", t, err)
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		log.Fatal().Msgf("could not find New function in %s: %s", t, err)
	}
	return newFunc.(func() interface{})().(Encoder)
}
