package convert

import "strconv"

type BoolConverter struct{}

func (c *BoolConverter) convert(v interface{}) (interface{}, error) {
	if v, ok := v.(string); ok {
		rst, err := strconv.ParseBool(v)
		if err != nil {
			return nil, err
		} else {
			return rst, err
		}
	}
	return nil, ErrConvertUnknownFormat
}
