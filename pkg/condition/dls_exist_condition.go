package condition

import "reflect"

type ExistCondition struct {
	paths []string
}

func NewExistCondition(paths []string) *ExistCondition {
	return &ExistCondition{paths}
}

func (c *ExistCondition) Pass(event map[string]interface{}) bool {
	var (
		o      map[string]interface{} = event
		length int                    = len(c.paths)
	)
	for _, path := range c.paths[:length-1] {
		if v, ok := o[path]; ok && v != nil {
			if reflect.TypeOf(v).Kind() == reflect.Map {
				o = v.(map[string]interface{})
			} else {
				return false
			}
		} else {
			return false
		}
	}

	if _, ok := o[c.paths[length-1]]; ok {
		return true
	}
	return false
}
