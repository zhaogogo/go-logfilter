package grok

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	a := `^(?P<logtime>\S+) (?P<name>\w+) (?P<status>\d+)$`
	p := regexp.MustCompile(`%{(\w+?)(?::(\w+?))?}`)
	res := p.FindAllStringSubmatch(a, -1)
	fmt.Printf("%q\n", res)
}
