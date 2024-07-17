package encoding

import (
	"github.com/zhaogogo/go-logfilter/encoding/json"
	"github.com/zhaogogo/go-logfilter/encoding/yaml"
	"strings"
)

func init() {
	RegisterCodec(yaml.Code{})
	RegisterCodec(json.Code{})
}

type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
}

var registeredCodecs = make(map[string]Codec)

func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}

func GetCodec(contentSubtype string) (codec Codec, ok bool) {
	codec, ok = registeredCodecs[contentSubtype]
	return
}
