package condition

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"reflect"
	"regexp"
	"strings"
)

type HasPrefixCondition struct {
	pat    *jsonpath.Compiled
	paths  []string
	prefix string
}

func NewHasPrefixCondition(c string) (*HasPrefixCondition, error) {
	if strings.HasPrefix(c, `HasPrefix($.`) {
		p := regexp.MustCompile(`^HasPrefix\((\$\..*),"(.*)"\)$`)
		r := p.FindStringSubmatch(c)
		if len(r) != 3 {
			return nil, fmt.Errorf("split jsonpath pattern/value error in `%s`", c)
		}

		value := r[2]
		pat, err := jsonpath.Compile(r[1])
		if err != nil {
			return nil, err
		}

		return &HasPrefixCondition{pat, nil, value}, nil
	}

	paths := make([]string, 0)
	c = strings.TrimSuffix(strings.TrimPrefix(c, "HasPrefix("), ")")
	for _, p := range strings.Split(c, ",") {
		paths = append(paths, strings.Trim(p, " "))
	}
	value := paths[len(paths)-1]
	paths = paths[:len(paths)-1]
	return &HasPrefixCondition{nil, paths, value}, nil
}

func (c *HasPrefixCondition) Pass(event map[string]interface{}) bool {
	if c.pat != nil {
		v, err := c.pat.Lookup(event)
		return err == nil && strings.HasPrefix(v.(string), c.prefix)
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
			return strings.HasPrefix(v.(string), c.prefix)
		}
	}
	return false
}
