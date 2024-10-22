package condition

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"reflect"
	"regexp"
	"strings"
)

type MatchCondition struct {
	pat    *jsonpath.Compiled
	paths  []string
	regexp *regexp.Regexp
}

func NewMatchCondition(c string) (*MatchCondition, error) {
	if strings.HasPrefix(c, `Match($.`) {
		p := regexp.MustCompile(`^Match\((\$\..*),"(.*)"\)$`)
		r := p.FindStringSubmatch(c)
		if len(r) != 3 {
			return nil, fmt.Errorf("split jsonpath pattern/value error in `%s`", c)
		}

		pat, err := jsonpath.Compile(r[1])
		if err != nil {
			return nil, err
		}

		value := r[2]
		regexp, err := regexp.Compile(value)
		if err != nil {
			return nil, err
		}

		return &MatchCondition{pat, nil, regexp}, nil
	}

	paths := make([]string, 0)
	c = strings.TrimSuffix(strings.TrimPrefix(c, "Match("), ")")
	for _, p := range strings.Split(c, ",") {
		paths = append(paths, strings.Trim(p, " "))
	}
	value := paths[len(paths)-1]
	paths = paths[:len(paths)-1]
	regexp, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}
	return &MatchCondition{nil, paths, regexp}, nil
}

func (c *MatchCondition) Pass(event map[string]interface{}) bool {
	if c.pat != nil {
		v, err := c.pat.Lookup(event)
		return err == nil && c.regexp.MatchString(v.(string))
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
			return c.regexp.MatchString(v.(string))
		}
	}
	return false
}
