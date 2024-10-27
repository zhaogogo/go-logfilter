package field

import "regexp"

var (
	matchp       = regexp.MustCompile(`^(\[.*?\])+$`)
	findp        = regexp.MustCompile(`(\[(.*?)\])`)
	matchGoTemp  = regexp.MustCompile(`{{.*}}`)
	matchESIndex = regexp.MustCompile(`%{.*?}`) //%{+YYYY.MM.dd}
	jsonPath     = regexp.MustCompile(`^\$\.`)

	FailedTagKey = "failed_tag"
)
