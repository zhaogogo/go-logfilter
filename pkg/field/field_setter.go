package field

import "github.com/zhaogogo/go-logfilter/pkg/field/set"

type SetField interface {
	SetField(event map[string]any, value any, overwrite bool) map[string]any
}

type FieldSetter interface {
	SetField(event map[string]any, value interface{}, overwrite bool) map[string]any
	Name() string
}

type FieldSet struct {
	name        string
	fieldsetter SetField
}

func (f *FieldSet) Name() string {
	return f.name
}

func NewFieldSet(template string, fieldsetter SetField) FieldSetter {
	return &FieldSet{
		name:        template,
		fieldsetter: fieldsetter,
	}
}

func (f *FieldSet) SetField(event map[string]any, value any, overwrite bool) map[string]any {
	if value == nil {
		return event
	}
	return f.fieldsetter.SetField(event, value, overwrite)
}

func NewFieldSetter(template string) FieldSetter {
	if matchp.Match([]byte(template)) {
		fields := make([]string, 0)
		for _, v := range findp.FindAllStringSubmatch(template, -1) {
			fields = append(fields, v[2])
		}
		if len(fields) == 1 {
			return NewFieldSet(template, set.NewOneLevelFieldSetter(fields[0]))
		}
		return NewFieldSet(template, set.NewMultiLevelFieldSetter(fields))
	} else {
		return NewFieldSet(template, set.NewOneLevelFieldSetter(template))
	}
}
