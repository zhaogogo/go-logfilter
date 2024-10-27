package field

import (
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
	render       ValueRender
	failedTag    bool
	failedTagKey string
}

func NewValueRender(template string, render ValueRender, failedTag bool) *ValueRende {
	return &ValueRende{
		name:         template,
		render:       render,
		failedTag:    failedTag,
		failedTagKey: FailedTagKey,
	}
}

func (vr *ValueRende) Render(event map[string]any) any {
	res := vr.render.Render(event)
	if res == nil {
		if vr.failedTag {
			v, ok := event[FailedTagKey]
			if ok {
				vv := v.([]any)
				vv = append(vv, vr.name)
				event[FailedTagKey] = vv
			} else {
				event[FailedTagKey] = []any{vr.name}
			}
		}
	}
	return res
}

// getValueRender matches all regexp pattern and return a ValueRender
// return nil if no pattern matched
func getValueRender(template string, failedTag bool) ValueRender {
	if matchp.Match([]byte(template)) {
		fields := make([]string, 0)
		for _, v := range findp.FindAllStringSubmatch(template, -1) {
			fields = append(fields, v[2])
		}

		if len(fields) == 1 {
			return NewValueRender(template, value.NewOneLevelValueRender(fields[0]), failedTag)
		}
		return NewValueRender(template, value.NewMultiLevelValueRender(fields), failedTag)
	}
	if matchGoTemp.Match([]byte(template)) {
		return NewValueRender(template, value.NewTemplateValueRender(template), failedTag)
	}
	if matchESIndex.Match([]byte(template)) {
		return NewValueRender(template, value.NewIndexRender(template), failedTag)

	}
	if jsonPath.Match([]byte(template)) {
		return NewValueRender(template, value.NewJsonpathRender(template), failedTag)
	}
	return nil
}

// GetValueRender return a ValueRender, and return LiteralValueRender if no pattern matched
func GetValueRender(template interface{}, failedTag bool) ValueRender {
	if temp, ok := template.(string); ok {
		r := getValueRender(temp, failedTag)
		if r != nil {
			return r
		}
	}
	return NewValueRender("LiteralValueRender", value.NewLiteralValueRender(template), failedTag)
}

// GetValueRender2 return a ValueRender, and return OneLevelValueRender("message") if no pattern matched
func GetValueRender2(template string, failedTag bool) ValueRender {
	r := getValueRender(template, failedTag)
	if r != nil {
		return r
	}
	return NewValueRender(template, value.NewOneLevelValueRender(template), failedTag)
}
