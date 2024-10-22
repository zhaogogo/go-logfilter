package condition

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"regexp"
	"strconv"
	"strings"
)

type INCondition struct {
	pat   *jsonpath.Compiled
	paths []string
	value interface{}
	fn    int
}

func NewINCondition(c string) (*INCondition, error) {
	var (
		pat   *jsonpath.Compiled
		paths []string
		value string
		err   error
	)

	if strings.HasPrefix(c, `IN($.`) {
		p := regexp.MustCompile(`^IN\((\$\..*),(.*)\)$`)
		r := p.FindStringSubmatch(c)
		if len(r) != 3 {
			return nil, fmt.Errorf("split jsonpath pattern/value error in `%s`", c)
		}

		if pat, err = jsonpath.Compile(r[1]); err != nil {
			return nil, err
		}

		value = r[2]
	} else {
		paths = make([]string, 0)
		c = strings.TrimSuffix(strings.TrimPrefix(c, "IN("), ")")
		for _, p := range strings.Split(c, ",") {
			paths = append(paths, strings.Trim(p, " "))
		}
		value = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	}

	if value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
		return &INCondition{pat, paths, value, len(paths)}, nil
	}
	if strings.Contains(value, ".") {
		if s, err := strconv.ParseFloat(value, 64); err == nil {
			return &INCondition{pat, paths, s, len(paths)}, nil
		}
		return nil, err
	}
	if s, err := strconv.ParseInt(value, 0, 32); err == nil {
		return &INCondition{pat, paths, int(s), len(paths)}, nil
	}
	return nil, err
}

func (c *INCondition) Pass(event map[string]interface{}) bool {
	if c.pat != nil {
		v, err := c.pat.Lookup(event)
		if err != nil {
			return false
		}
		if l, ok := v.([]interface{}); !ok {
			return false
		} else {
			for _, e := range l {
				if c.value == e {
					return true
				}
			}
		}
		return false
	}

	var (
		o map[string]interface{} = event
	)

	for _, path := range c.paths[:c.fn-1] {
		if v, ok := o[path]; ok && v != nil {
			if o, ok = v.(map[string]interface{}); !ok {
				return false
			}
		} else {
			return false
		}
	}

	if v, ok := o[c.paths[c.fn-1]]; ok {
		if l, ok := v.([]interface{}); !ok {
			return false
		} else {
			for _, e := range l {
				if c.value == e {
					return true
				}
			}
		}
	}
	return false

}
