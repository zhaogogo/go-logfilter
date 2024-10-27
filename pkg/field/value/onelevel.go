package valuerender

type OneLevelValueRender struct {
	field string
}

func NewOneLevelValueRender(template string) *OneLevelValueRender {
	return &OneLevelValueRender{
		field: template,
	}
}

func (l *OneLevelValueRender) Render(event map[string]interface{}) interface{} {
	if value, ok := event[l.field]; ok {
		return value
	}
	return nil
}
