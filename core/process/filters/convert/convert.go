package convert

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/field"
)

type Converter interface {
	convert(v interface{}) (interface{}, error)
}

var ErrConvertUnknownFormat error = errors.New("unknown format")

type ConveterAndRender struct {
	converter    Converter
	valueRender  field.ValueRender
	removeIfFail bool
	settoIfFail  interface{}
	settoIfNil   interface{}
}

type ConvertFilter struct {
	config map[string]interface{}
	fields map[field.FieldSetter]ConveterAndRender
	//failed_tag bool
}

func New(config map[string]interface{}) topology.Filter {
	plugin := &ConvertFilter{
		config: config,
		fields: make(map[field.FieldSetter]ConveterAndRender),
		//failed_tag: config["failed_tag"].(bool),
	}

	if fieldsValue, ok := config["fields"]; ok {
		for f, vI := range fieldsValue.(map[string]interface{}) {
			v := vI.(map[string]interface{})
			overwrite := false
			fieldSetter := field.NewFieldSetter("convert", f, &overwrite)
			if fieldSetter == nil {
				log.Fatal().Msgf("could build field setter from %#v", f)
			}

			to := v["to"].(string)
			remove_if_fail := false
			if I, ok := v["remove_if_fail"]; ok {
				remove_if_fail = I.(bool)
			}
			setto_if_fail := v["setto_if_fail"]
			setto_if_nil := v["setto_if_nil"]

			var converter Converter
			if to == "float" {
				converter = &FloatConverter{}
			} else if to == "int" {
				converter = &IntConverter{}
			} else if to == "uint" {
				converter = &UIntConverter{}
			} else if to == "bool" {
				converter = &BoolConverter{}
			} else if to == "string" {
				converter = &StringConverter{}
			} else if to == "array(int)" {
				converter = &ArrayIntConverter{}
			} else if to == "array(float)" {
				converter = &ArrayFloatConverter{}
			} else {
				log.Fatal().Msg("can only convert to int/float/bool/array(int)/array(float)")
			}

			plugin.fields[fieldSetter] = ConveterAndRender{
				converter:    converter,
				valueRender:  field.GetValueRender2("convert", f, true),
				removeIfFail: remove_if_fail,
				settoIfFail:  setto_if_fail,
				settoIfNil:   setto_if_nil,
			}
		}
	} else {
		log.Fatal().Msg("fileds must be set in convert filter plugin")
	}
	return plugin
}

func (plugin *ConvertFilter) Filter(event map[string]any) (map[string]any, error) {
	for fs, conveterAndRender := range plugin.fields {
		originanV := conveterAndRender.valueRender.Render(event)
		if originanV == nil {
			//if plugin.failed_tag {
			//	event = field.SetFailedTags(event, field.FailedTagKey, fmt.Sprintf("%s field convert value render got nil", fs.Name()))
			//}
			if conveterAndRender.settoIfNil != nil {
				event = fs.SetField(event, conveterAndRender.settoIfNil, true)
			}
			continue
		}
		v, err := conveterAndRender.converter.convert(originanV)
		if err == nil {
			event = fs.SetField(event, v, true)
		} else {
			log.Info().Msgf("convert 字段%s, error: %s", fs.Name(), err)
			//if plugin.failed_tag {
			//	event = field.SetFailedTags(event, field.FailedTagKey, fmt.Sprintf("%s field convert failed", fs.Name()))
			//}
			if conveterAndRender.removeIfFail {
				event = fs.SetField(event, nil, true)
			} else if conveterAndRender.settoIfFail != nil {
				event = fs.SetField(event, conveterAndRender.settoIfFail, true)
			}
		}
	}
	return event, nil
}
