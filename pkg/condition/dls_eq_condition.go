package condition

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"regexp"
	"strconv"
	"strings"
)

type EQCondition struct {
	pat   *jsonpath.Compiled
	paths []string
	value interface{}
	fn    int
}

func NewEQCondition(c string) (*EQCondition, error) {
	var (
		pat   *jsonpath.Compiled
		paths []string
		value string
		err   error
	)

	if strings.HasPrefix(c, `EQ($.`) {
		p := regexp.MustCompile(`^EQ\((\$\..*),(.*)\)$`)
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
		c = strings.TrimSuffix(strings.TrimPrefix(c, "EQ("), ")")
		for _, p := range strings.Split(c, ",") {
			paths = append(paths, strings.Trim(p, " "))
		}
		value = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	}

	if value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
		return &EQCondition{pat, paths, value, len(paths)}, nil
	}

	if value == "nil" {
		return &EQCondition{pat, paths, nil, len(paths)}, nil
	}

	if strings.Contains(value, ".") {
		if s, err := strconv.ParseFloat(value, 64); err == nil {
			return &EQCondition{pat, paths, s, len(paths)}, nil
		}
		return nil, err
	}
	if s, err := strconv.ParseInt(value, 0, 32); err == nil {
		return &EQCondition{pat, paths, s, len(paths)}, nil
	} else {
		return nil, err
	}
}

func (c *EQCondition) Pass(event map[string]interface{}) bool {
	if c.pat != nil {
		v, err := c.pat.Lookup(event)
		return err == nil && equal(v, c.value)
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
		return equal(v, c.value)
	}
	return false
}

func equal(src, target interface{}) bool {
	if n, ok := src.(json.Number); ok {
		if tValue, ok := target.(int64); ok {
			if intV, err := n.Int64(); err == nil {
				return intV == tValue
			}
			return false
		}
		if tValue, ok := target.(float64); ok {
			if floatV, err := n.Float64(); err == nil {
				return floatV == tValue
			}
			return false
		}
	}
	return src == target
}
