package field

import (
	"github.com/rs/zerolog/log"
	value "github.com/zhaogogo/go-logfilter/pkg/field/value"
)

var (
	_ ValueRender = &value.OneLevelValueRender{}
	_ ValueRender = &value.MultiLevelValueRender{}
	_ ValueRender = &value.JsonpathRender{}
	_ ValueRender = &value.TemplateValueRender{}
)

type ValueRender interface {
	Render(map[string]interface{}) interface{}
}

type ValueRende struct {
	name         string
	template     string
	render       ValueRender
	failedTag    *bool
	failedTagKey string
}

func NewValueRender(name string, template string, render ValueRender, failedTag *bool) *ValueRende {
	return &ValueRende{
		name:         name,
		template:     template,
		render:       render,
		failedTag:    failedTag,
		failedTagKey: FailedTagKey,
	}
}

func (vr *ValueRende) Render(event map[string]any) any {
	res := vr.render.Render(event)
	if res == nil {
		log.Debug().Msgf("%s %s value render failed", vr.name, vr.template)
		if vr.failedTag != nil && *vr.failedTag {
			v, ok := event[FailedTagKey]
			if ok {
				vv := v.([]any)
				vv = append(vv, vr.template)
				event[FailedTagKey] = vv
			} else {
				event[FailedTagKey] = []any{vr.template}
			}
		}
	}
	return res
}

// getValueRender matches all regexp pattern and return a ValueRender
// return nil if no pattern matched
func getValueRender(name string, template string, failedTag *bool) ValueRender {
	if matchp.Match([]byte(template)) {
		fields := make([]string, 0)
		for _, v := range findp.FindAllStringSubmatch(template, -1) {
			fields = append(fields, v[2])
		}

		if len(fields) == 1 {
			return NewValueRender(name, template, value.NewOneLevelValueRender(fields[0]), failedTag)
		}
		return NewValueRender(name, template, value.NewMultiLevelValueRender(fields), failedTag)
	}
	if matchGoTemp.Match([]byte(template)) {
		return NewValueRender(name, template, value.NewTemplateValueRender(template), failedTag)
	}
	if matchESIndex.Match([]byte(template)) {
		return NewValueRender(name, template, value.NewIndexRender(template), failedTag)

	}
	if jsonPath.Match([]byte(template)) {
		return NewValueRender(name, template, value.NewJsonpathRender(template), failedTag)
	}
	return nil
}

// GetValueRender return a ValueRender, and return LiteralValueRender if no pattern matched
func GetValueRender(name string, template interface{}, failedTag *bool) ValueRender {
	if temp, ok := template.(string); ok {
		r := getValueRender(name, temp, failedTag)
		if r != nil {
			return r
		}
	}
	return NewValueRender(name, "LiteralValueRender", value.NewLiteralValueRender(template), failedTag)
}

// GetValueRender2 return a ValueRender, and return OneLevelValueRender("message") if no pattern matched
func GetValueRender2(name string, template string, failedTag *bool) ValueRender {
	r := getValueRender(name, template, failedTag)
	if r != nil {
		return r
	}
	return NewValueRender(name, template, value.NewOneLevelValueRender(template), failedTag)
}
