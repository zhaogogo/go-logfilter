package condition

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type OPNode struct {
	//name      string
	op        int
	left      *OPNode
	right     *OPNode
	condition Condition //leaf node has condition
	pos       int
}

func (root *OPNode) Pass(event map[string]interface{}) bool {
	if root.condition != nil {
		return root.condition.Pass(event)
	}

	if root.op == _op_and {
		return root.left.Pass(event) && root.right.Pass(event)
	}
	if root.op == _op_or {
		return root.left.Pass(event) || root.right.Pass(event)
	}
	if root.op == _op_not {
		return !root.right.Pass(event)
	}
	return false
}

func NewSingleCondition(c string) (Condition, error) {
	original_c := c

	// Exist
	if matched, _ := regexp.MatchString(`^Exist\(.*\)$`, c); matched {
		c = strings.TrimSuffix(strings.TrimPrefix(c, "Exist("), ")")
		paths := make([]string, 0)
		for _, p := range strings.Split(c, ",") {
			paths = append(paths, strings.Trim(p, " "))
		}
		return NewExistCondition(paths), nil
	}

	// IN
	if matched, _ := regexp.MatchString(`^IN\(.*\)$`, c); matched {
		return NewINCondition(c)
	}

	// EQ
	if matched, _ := regexp.MatchString(`^EQ\(.*\)$`, c); matched {
		return NewEQCondition(c)
	}

	// HasPrefix
	if matched, _ := regexp.MatchString(`^HasPrefix\(.*\)$`, c); matched {
		return NewHasPrefixCondition(c)
	}

	// HasSuffix
	if matched, _ := regexp.MatchString(`^HasSuffix\(.*\)$`, c); matched {
		return NewHasSuffixCondition(c)
	}

	// Contains
	if matched, _ := regexp.MatchString(`^Contains\(.*\)$`, c); matched {
		return NewContainsCondition(c)
	}

	// ContainsAny
	if matched, _ := regexp.MatchString(`^ContainsAny\(.*\)$`, c); matched {
		paths := make([]string, 0)
		c = strings.TrimSuffix(strings.TrimPrefix(c, "ContainsAny("), ")")
		for _, p := range strings.Split(c, ",") {
			paths = append(paths, strings.Trim(p, " "))
		}
		value := paths[len(paths)-1]
		paths = paths[:len(paths)-1]
		return NewContainsAnyCondition(paths, value), nil
	}

	// Match
	if matched, _ := regexp.MatchString(`^Match\(.*\)$`, c); matched {
		return NewMatchCondition(c)
	}

	// Random
	if matched, _ := regexp.MatchString(`^Random\(.*\)$`, c); matched {
		c = strings.TrimSuffix(strings.TrimPrefix(c, "Random("), ")")
		if value, err := strconv.ParseInt(c, 0, 32); err != nil {
			return nil, err
		} else {
			return NewRandomCondition(int(value)), nil
		}
	}

	// Before
	if matched, _ := regexp.MatchString(`^Before\(.*\)$`, c); matched {
		c = strings.TrimSuffix(strings.TrimPrefix(c, "Before("), ")")
		return NewBeforeCondition(c), nil
	}

	// After
	if matched, _ := regexp.MatchString(`^After\(.*\)$`, c); matched {
		c = strings.TrimSuffix(strings.TrimPrefix(c, "After("), ")")
		return NewAfterCondition(c), nil
	}

	return nil, fmt.Errorf("could not build Condition from `%s`", original_c)
}
