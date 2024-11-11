package json

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Jsoner struct {
	TimestampKey string
	UseNumber    bool
}

func (jd *Jsoner) Decode(value []byte) map[string]interface{} {
	rst := make(map[string]interface{})
	d := json.NewDecoder(bytes.NewReader(value))

	if jd.UseNumber {
		d.UseNumber()
	}
	err := d.Decode(&rst)
	if err != nil || d.More() {
		return map[string]interface{}{
			jd.TimestampKey: time.Now().Format("2006-01-02T15:04:05.000 -0700"),
			"message":       string(value),
		}
	}
	if _, ok := rst[jd.TimestampKey]; !ok {
		rst[jd.TimestampKey] = time.Now()
	}
	return rst
}

func (j *Jsoner) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
