package condition

import (
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

func NewCondition(name string, template string) Condition {
	original_c := template

	template = strings.Trim(template, " ")

	if matched, _ := regexp.MatchString(`^{{.*}}$`, template); matched {
		return NewTemplateCondition(name, template)
	}

	if root, err := parseBoolTree(template); err != nil {
		log.Panic().Msgf("could not build Condition from `%s` : %s", original_c, err)
		return nil
	} else {
		return root
	}

}
