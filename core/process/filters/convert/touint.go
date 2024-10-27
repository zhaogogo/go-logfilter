package convert

import (
	"github.com/spf13/cast"
)

type UIntConverter struct{}

func (c *UIntConverter) convert(v interface{}) (interface{}, error) {
	return cast.ToUint64E(v)
}
