package convert

import "encoding/json"

type ArrayIntConverter struct{}

func (c *ArrayIntConverter) convert(v interface{}) (interface{}, error) {
	if v1, ok1 := v.([]interface{}); ok1 {
		var t2 = []int{}
		for _, i := range v1 {
			j, err := i.(json.Number).Int64()
			// j, err := strconv.ParseInt(i.String(), 0, 64)
			if err != nil {
				return nil, ErrConvertUnknownFormat
			}
			t2 = append(t2, (int)(j))
		}
		return t2, nil
	}
	return nil, ErrConvertUnknownFormat
}
