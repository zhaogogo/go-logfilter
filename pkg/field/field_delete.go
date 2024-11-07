package field

import (
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/pkg/field/del"
	"regexp"
	"strings"
)

var (
	_ FieldDelete = &del.OneLevelFieldDeleter{}
	_ FieldDelete = &del.MultiLevelFieldDeleter{}
)

type FieldDelete interface {
	Delete(map[string]any)
}

func NewFieldDeleter(c map[string]any, pluginname string) []FieldDelete {

	delconfs, ok := c["delete_fields"].([]any)
	if !ok {
		log.Debug().Msgf("%s get delete_fields assets failed, got %T ignore", pluginname, c["delete_fields"])
		return nil
	}
	res := make([]FieldDelete, 0, len(delconfs))
	for _, delc := range delconfs {
		if template, ok := delc.(string); ok {
			res = append(res, NewFieldDel(template))
		} else {
			log.Debug().Msgf("%s 获取delete字段失败, got: %T, want: string ignore", pluginname, delc)
		}
	}
	return res
}

func NewFieldDel(template string) FieldDelete {
	matchp, _ := regexp.Compile(`(\[.*?\])+`)
	findp, _ := regexp.Compile(`(\[(.*?)\])`)
	if matchp.Match([]byte(template)) {
		fields := make([]string, 0)
		for _, v := range findp.FindAllStringSubmatch(template, -1) {
			if v[2] != "" {
				fields = append(fields, strings.TrimSpace(v[2]))
			}
		}
		if len(fields) == 0 {
			return del.NewOneLevelFieldDeleter(template)
		}
		return del.NewMultiLevelFieldDeleter(fields)
	} else {
		return del.NewOneLevelFieldDeleter(template)
	}
}
