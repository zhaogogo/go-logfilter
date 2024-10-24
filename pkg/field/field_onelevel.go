package field

type OneLevelFieldSetter struct {
	field string
}

func NewOneLevelFieldSetter(field string) *OneLevelFieldSetter {
	r := &OneLevelFieldSetter{
		field: field,
	}
	return r
}

func (fs *OneLevelFieldSetter) SetField(event map[string]interface{}, value interface{}, overwrite bool) map[string]interface{} {
	if _, ok := event[fs.field]; !ok || overwrite {
		event[fs.field] = value
	}
	return event
}
