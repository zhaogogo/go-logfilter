package valuerender

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	a := "{x1}{x2}"
	//a := "%{+YYYY.MM.dd}{xx}"
	r, _ := regexp.Compile(`{(.*?)}`)
	res := r.FindAll([]byte(a), -1)
	fmt.Printf("%q\n", res)
}
