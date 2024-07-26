package textcodec

import (
	"github.com/zhaogogo/go-logfilter/textcodec/json"
	"k8s.io/klog/v2"
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
		klog.Fatalf("could not open %s: %s", t, err)
	}
	newFunc, err := p.Lookup("New")
	if err != nil {
		klog.Fatalf("could not find New function in %s: %s", t, err)
	}
	return newFunc.(func() interface{})().(Encoder)
}
