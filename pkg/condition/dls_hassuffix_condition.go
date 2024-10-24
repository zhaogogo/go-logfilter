package condition

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"reflect"
	"regexp"
	"strings"
)

type HasSuffixCondition struct {
	pat    *jsonpath.Compiled
	paths  []string
	suffix string
}

func NewHasSuffixCondition(c string) (*HasSuffixCondition, error) {
	if strings.HasPrefix(c, `HasSuffix($.`) {
		p := regexp.MustCompile(`^HasSuffix\((\$\..*),"(.*)"\)$`)
		r := p.FindStringSubmatch(c)
		if len(r) != 3 {
			return nil, fmt.Errorf("split jsonpath pattern/value error in `%s`", c)
		}

		value := r[2]
		pat, err := jsonpath.Compile(r[1])
		if err != nil {
			return nil, err
		}

		return &HasSuffixCondition{pat, nil, value}, nil
	}

	paths := make([]string, 0)
	c = strings.TrimSuffix(strings.TrimPrefix(c, "HasSuffix("), ")")
	for _, p := range strings.Split(c, ",") {
		paths = append(paths, strings.Trim(p, " "))
	}
	value := paths[len(paths)-1]
	paths = paths[:len(paths)-1]
	return &HasSuffixCondition{nil, paths, value}, nil
}

func (c *HasSuffixCondition) Pass(event map[string]interface{}) bool {
	if c.pat != nil {
		v, err := c.pat.Lookup(event)
		return err == nil && strings.HasSuffix(v.(string), c.suffix)
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
			return strings.HasSuffix(v.(string), c.suffix)
		}
	}
	return false
}
