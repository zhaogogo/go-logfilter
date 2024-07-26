package textcodec

import (
	"github.com/zhaogogo/go-logfilter/textcodec/json"
	"github.com/zhaogogo/go-logfilter/textcodec/plain"
	"k8s.io/klog/v2"
	"plugin"
)

type Decoder interface {
	Decode([]byte) map[string]interface{}
}

func NewDecoder(t string) Decoder {
	switch t {
	case "plain":
		return &plain.PlainDecoder{}
	case "json":
		return &json.Jsoner{UseNumber: true}
	case "json:not_usenumber":
		return &json.Jsoner{UseNumber: false}
	default:
		p, err := plugin.Open(t)
		if err != nil {
			klog.Fatalf("could not open %s: %s", t, err)
		}
		newFunc, err := p.Lookup("New")
		if err != nil {
			klog.Fatalf("could not find New function in %s: %s", t, err)
		}
		return newFunc.(func() interface{})().(Decoder)
	}
}
