package condition

import (
	"reflect"
	"strings"
)

type ContainsAnyCondition struct {
	paths     []string
	substring string
}

func NewContainsAnyCondition(paths []string, substring string) *ContainsAnyCondition {
	return &ContainsAnyCondition{paths, substring}
}

func (c *ContainsAnyCondition) Pass(event map[string]interface{}) bool {
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

	if v, ok := o[c.paths[length-1]]; ok && v != nil {
		if reflect.TypeOf(v).Kind() == reflect.String {
			return strings.ContainsAny(v.(string), c.substring)
		}
	}
	return false
}
