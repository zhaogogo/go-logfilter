package field

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

func NewAddFields(name string, c map[string]interface{}) []map[FieldSetter]ValueRender {
	addFieldsRes := []map[FieldSetter]ValueRender{}
	var (
		failedTag bool  = true
		overwrite *bool = nil
	)
	add_fields, ok := c["add_fields"]
	if !ok {
		log.Warn().Msgf("%s add_fields field not exists", name)
		return nil
	}
	if add_fields == nil {
		log.Warn().Msgf("%s add_fields field value is nil", name)
		return nil
	}
	if add_fields == nil {
		log.Warn().Msgf("%s add_fields field value is nil", name)
		return nil
	}

	add_fieldsConf, ok := add_fields.(map[string]any)
	if !ok {
		log.Warn().Msgf("add_fields conf asset failed, got=%T", add_fields)
		return nil
	}

	failedTagAny, ok := add_fieldsConf["failed_tag"]
	if !ok {
		log.Warn().Msgf("%s add_fields add_fieldsConf failed_tag field not exists", name)
	} else {
		failedTagb, ok := failedTagAny.(bool)
		if ok {
			failedTag = failedTagb
		}
	}
	overwriteAny, ok := add_fieldsConf["overwrite"]
	if !ok {
		log.Warn().Msgf("%s add_fields add_fieldsConf overwriteAny field not exists", name)
	} else {
		overwriteb, ok := overwriteAny.(bool)
		if ok {
			overwrite = &overwriteb
		}
	}

	add_fields_kv, ok := add_fieldsConf["kv"]
	if !ok {
		log.Warn().Msgf("add_fields conf kv field not exists ignore")
		return nil
	}
	add_fields_kvs, ok := add_fields_kv.([]any)
	if !ok {
		log.Warn().Msgf("add_fields conf kv field asset failed, got=%T", add_fields_kv)
		return nil
	}

	for _, add_fields := range add_fields_kvs {
		vv, ok := add_fields.(map[string]interface{})
		if !ok {
			log.Warn().Msgf("add_fields conf kv field asset failed, got=%T", add_fields)
			return nil
		}

		addFields := make(map[FieldSetter]ValueRender)
		for template, v := range vv {
			fieldSetter := NewFieldSetter(fmt.Sprintf("%s add_field", name), template, overwrite)
			addFields[fieldSetter] = GetValueRender(fmt.Sprintf("%s add_field", name), v, &failedTag)
		}
		addFieldsRes = append(addFieldsRes, addFields)
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
