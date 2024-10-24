package filters

import (
	"github.com/zhaogogo/go-logfilter/core"
)

func init() {
	Register("Grok", newGrokFilter)
}

func newGrokFilter(m map[string]interface{}) core.Processer {

}
