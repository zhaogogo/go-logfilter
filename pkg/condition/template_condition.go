package condition

import "github.com/zhaogogo/go-logfilter/pkg/field"

var _ Condition = &TemplateCondition{}

func NewTemplateCondition(name string, template string) *TemplateCondition {
	failedTag := false
	return &TemplateCondition{
		ifCondition: field.GetValueRender(name, template, &failedTag),
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
