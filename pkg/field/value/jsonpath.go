package valuerender

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
)

type JsonpathRender struct {
	Pat *jsonpath.Compiled
}

func NewJsonpathRender(template string) *JsonpathRender {
	pat, err := jsonpath.Compile(template)
	if err != nil {
		panic(fmt.Sprintf("json path compile `%s` error: %s", template, err))
	}
	return &JsonpathRender{
		Pat: pat,
	}
}

func (l *JsonpathRender) Render(event map[string]interface{}) interface{} {
	if value, ok := l.Pat.Lookup(event); ok == nil {
		return value
	}
	return nil
}
