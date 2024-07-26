package json

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Jsoner struct {
	UseNumber bool
}

func (jd *Jsoner) Decode(value []byte) map[string]interface{} {
	rst := make(map[string]interface{})
	rst["@timestamp"] = time.Now()
	d := json.NewDecoder(bytes.NewReader(value))

	if jd.UseNumber {
		d.UseNumber()
	}
	err := d.Decode(&rst)
	if err != nil || d.More() {
		return map[string]interface{}{
			"@timestamp": time.Now(),
			"message":    string(value),
		}
	}

	return rst
}

func (j *Jsoner) Encode(v interface{}) ([]byte, error) {

	return json.Marshal(v)
}
