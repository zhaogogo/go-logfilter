package condition

import "github.com/zhaogogo/go-logfilter/field"

var _ Condition = &TemplateCondition{}

func NewTemplateCondition(c string) *TemplateCondition {
	return &TemplateCondition{
		ifCondition: field.GetValueRender(c),
		ifResult:    "y",
	}
}

type TemplateCondition struct {
	ifCondition field.ValueRender
	ifResult    string
}

func (t *TemplateCondition) Pass(event map[string]interface{}) bool {
	r := t.ifCondition.Render(event)
	if r == nil || r.(string) != t.ifResult {
		return false
	}
	return true
}
