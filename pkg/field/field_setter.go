package field

type FieldSetter interface {
	SetField(event map[string]interface{}, value interface{}, overwrite bool) map[string]interface{}
}

func NewFieldSetter(template string) FieldSetter {
	if matchp.Match([]byte(template)) {
		fields := make([]string, 0)
		for _, v := range findp.FindAllStringSubmatch(template, -1) {
			fields = append(fields, v[2])
		}
		if len(fields) == 1 {
			return NewOneLevelFieldSetter(fields[0])
		}
		return NewMultiLevelFieldSetter(fields)
	} else {
		return NewOneLevelFieldSetter(template)
	}
}
