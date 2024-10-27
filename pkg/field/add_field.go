package field

import (
	"github.com/rs/zerolog/log"
)

func NewAddFields(c map[string]interface{}, failedTag bool) []map[FieldSetter]ValueRender {
	addFieldsRes := []map[FieldSetter]ValueRender{}
	if add_fields, ok := c["add_fields"]; ok && add_fields != nil {
		add_fieldss, ok := add_fields.([]any)
		if !ok {
			log.Warn().Msgf("add_fields asset failed, %T", add_fieldss)
			return nil
		}
		//fmt.Printf("--->%T\n", add_fieldss)
		for _, add_fields := range add_fieldss {
			//fmt.Printf("--->%T %#v\n", add_fields, add_fields)
			vv, ok := add_fields.(map[string]interface{})
			if !ok {
				log.Warn().Msgf("add_fields asset failed, %T", add_fields)
				return nil
			}

			addFields := make(map[FieldSetter]ValueRender)
			for k, v := range vv {
				fieldSetter := NewFieldSetter(k)
				addFields[fieldSetter] = GetValueRender(v, failedTag)
			}
			//return addFields
			//fmt.Println("xxxx", len(addFields))
			addFieldsRes = append(addFieldsRes, addFields)
		}
	}
	return addFieldsRes
}

func SetFailedTags(event map[string]any, failedTagKey string, failedValue any) map[string]any {
	v, ok := event[failedTagKey]
	if ok {
		vv := v.([]any)
		vv = append(vv, failedValue)
		event[failedTagKey] = vv
	} else {
		event[failedTagKey] = []any{failedValue}
	}
	return event
}
