package field

import (
	"github.com/rs/zerolog/log"
)

func NewAddFields(c map[string]interface{}) map[FieldSetter]ValueRender {
	if add_fields, ok := c["add_fields"]; ok && add_fields != nil {
		vv, ok := add_fields.(map[string]interface{})
		if !ok {
			log.Warn().Msgf("add_fields asset failed, %T", add_fields)
			return nil
		}
		addFields := make(map[FieldSetter]ValueRender)
		for k, v := range vv {
			fieldSetter := NewFieldSetter(k)
			addFields[fieldSetter] = GetValueRender(v)
		}
		return addFields
	}
	return nil
}
