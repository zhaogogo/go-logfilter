package convert

import "github.com/spf13/cast"

type FloatConverter struct{}

func (c *FloatConverter) convert(v interface{}) (interface{}, error) {
	return cast.ToFloat64E(v)
}
