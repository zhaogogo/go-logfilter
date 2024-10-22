package condition

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"reflect"
	"regexp"
	"strings"
)

type ContainsCondition struct {
	pat       *jsonpath.Compiled
	paths     []string
	substring string
}

func NewContainsCondition(c string) (*ContainsCondition, error) {
	if strings.HasPrefix(c, `Contains($.`) {
		p := regexp.MustCompile(`^Contains\((\$\..*),"(.*)"\)$`)
		r := p.FindStringSubmatch(c)
		if len(r) != 3 {
			return nil, fmt.Errorf("split jsonpath pattern/value error in `%s`", c)
		}

		value := r[2]
		pat, err := jsonpath.Compile(r[1])
		if err != nil {
			return nil, err
		}

		return &ContainsCondition{pat, nil, value}, nil
	}
	paths := make([]string, 0)
	c = strings.TrimSuffix(strings.TrimPrefix(c, "Contains("), ")")
	for _, p := range strings.Split(c, ",") {
		paths = append(paths, strings.Trim(p, " "))
	}
	value := paths[len(paths)-1]
	paths = paths[:len(paths)-1]
	return &ContainsCondition{nil, paths, value}, nil
}

func (c *ContainsCondition) Pass(event map[string]interface{}) bool {
	if c.pat != nil {
		v, err := c.pat.Lookup(event)
		return err == nil && strings.Contains(v.(string), c.substring)
	}

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
			return strings.Contains(v.(string), c.substring)
		}
	}
	return false
}
