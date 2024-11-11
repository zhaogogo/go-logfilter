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
	template    string
	fieldsetter SetField
	overwrite   *bool
}

func (f *FieldSet) Name() string {
	return f.name
}

func newFieldSet(name string, template string, fieldsetter SetField, overwrite *bool) FieldSetter {
	return &FieldSet{
		name:        name,
		template:    template,
		fieldsetter: fieldsetter,
		overwrite:   overwrite,
	}
}

func (f *FieldSet) SetField(event map[string]any, value any, overwrite bool) map[string]any {
	if value == nil {
		return event
	}

	if f.overwrite != nil {
		overwrite = *f.overwrite
	}

	return f.fieldsetter.SetField(event, value, overwrite)
}

func NewFieldSetter(name string, template string, overwrite *bool) FieldSetter {
	if matchp.Match([]byte(template)) {
		fields := make([]string, 0)
		for _, v := range findp.FindAllStringSubmatch(template, -1) {
			fields = append(fields, v[2])
		}
		if len(fields) == 1 {
			return newFieldSet(name, template, set.NewOneLevelFieldSetter(fields[0]), overwrite)
		}
		return newFieldSet(name, template, set.NewMultiLevelFieldSetter(fields), overwrite)
	} else {
		return newFieldSet(name, template, set.NewOneLevelFieldSetter(template), overwrite)
	}
}
