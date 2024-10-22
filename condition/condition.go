package condition

import (
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

func NewCondition(c string) Condition {
	original_c := c

	c = strings.Trim(c, " ")

	if matched, _ := regexp.MatchString(`^{{.*}}$`, c); matched {
		return NewTemplateCondition(c)
	}

	if root, err := parseBoolTree(c); err != nil {
		log.Panic().Msgf("could not build Condition from `%s` : %s", original_c, err)
		return nil
	} else {
		return root
	}

}
