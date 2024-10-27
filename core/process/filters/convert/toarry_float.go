package convert

import "encoding/json"

type ArrayFloatConverter struct{}

func (c *ArrayFloatConverter) convert(v interface{}) (interface{}, error) {
	if v1, ok1 := v.([]interface{}); ok1 {
		var t2 = []float64{}
		for _, i := range v1 {
			j, err := i.(json.Number).Float64()
			if err != nil {
				return nil, ErrConvertUnknownFormat
			}
			t2 = append(t2, (float64)(j))
		}
		return t2, nil
	}
	return nil, ErrConvertUnknownFormat
}
