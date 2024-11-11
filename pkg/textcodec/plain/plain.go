package plain

import "time"

type PlainDecoder struct {
	TimestampKey string
}

func (d *PlainDecoder) Decode(value []byte) map[string]interface{} {
	rst := make(map[string]interface{})
	rst[d.TimestampKey] = time.Now().Format("2006-01-02T15:04:05.000 -0700")
	rst["message"] = string(value)
	return rst
}
