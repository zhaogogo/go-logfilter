package field

type LiteralValueRender struct {
	value interface{}
}

func NewLiteralValueRender(template interface{}) *LiteralValueRender {
	return &LiteralValueRender{template}
}

func (r *LiteralValueRender) Render(event map[string]interface{}) interface{} {
	return r.value
}
