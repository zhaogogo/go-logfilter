package plain

import "time"

type PlainDecoder struct {
	TimestampKey string
}

func (d *PlainDecoder) Decode(value []byte) map[string]interface{} {
	rst := make(map[string]interface{})
	rst[d.TimestampKey] = time.Now()
	rst["message"] = string(value)
	return rst
}
