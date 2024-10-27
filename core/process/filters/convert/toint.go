package convert

import "github.com/spf13/cast"

type IntConverter struct{}

func (c *IntConverter) convert(v interface{}) (interface{}, error) {
	return cast.ToInt64E(v)
}
