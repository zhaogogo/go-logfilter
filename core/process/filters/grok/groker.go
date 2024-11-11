package grok

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/core/topology"
	"github.com/zhaogogo/go-logfilter/pkg/field"
)

//
//func init() {
//	Register("Grok", newGrokFilter)
//}

type GrokFilter struct {
	config    map[string]any
	overwrite bool
	groks     []*Grok
	target    string
	src       string
	vr        field.ValueRender
}

func New(c map[string]any) topology.Filter {
	var patternPaths []string = make([]string, 0, 10)
	if pattern_paths, ok := c["pattern_paths"]; ok {
		for _, p := range pattern_paths.([]interface{}) {
			patternPaths = append(patternPaths, p.(string))
		}
	}
	ignoreBlank := true
	if i, ok := c["ignore_blank"]; ok {
		ignoreBlank = i.(bool)
	}
	groks := make([]*Grok, 0)
	if matchValue, ok := c["match"]; ok {
		match := matchValue.([]interface{})
		for _, mValue := range match {
			groks = append(groks, NewGrok(mValue.(string), patternPaths, ignoreBlank))
		}
	} else {
		log.Fatal().Msgf("grok filter plugin match field must be slience, got %T", c["match"])
	}
	gf := &GrokFilter{
		config:    c,
		groks:     groks,
		overwrite: false,
		target:    "",
	}
	var valueRenderOverwrite *bool = nil
	if overwrite, ok := c["overwrite"]; ok {
		if overwrite, ok := overwrite.(bool); ok {
			gf.overwrite = overwrite
			valueRenderOverwrite = &overwrite
		}
	}

	if srcValue, ok := c["src"]; ok {
		gf.src = srcValue.(string)
	} else {
		gf.src = "message"
	}
	gf.vr = field.GetValueRender2("gork", gf.src, valueRenderOverwrite)

	if target, ok := c["target"]; ok {
		gf.target = target.(string)
	}

	return gf
}

func (g *GrokFilter) Filter(event map[string]interface{}) (map[string]interface{}, error) {
	input := g.vr.Render(event)
	if input == nil {
		log.Error().Msgf("grok filter plugin field render value failed, event=%#v", event)
		return event, errors.New("grok filter plugin field render value failed")
	}
	i, ok := input.(string)
	if !ok {
		return event, errors.New("grok filter plugin field value result is not string")
	}
	for _, grok := range g.groks {
		rst := grok.grok(i)
		if len(rst) == 0 {
			continue
		}
		if g.target == "" {
			if g.overwrite {
				for fie, val := range rst {
					event[fie] = val
				}
			} else {
				for fie, val := range rst {
					if _, exists := event[fie]; !exists {
						event[fie] = val
					}
				}
			}

		} else {
			target := make(map[string]string)
			for fie, val := range rst {
				target[fie] = val
			}

			if g.overwrite {
				event[g.target] = target
			} else {
				if _, exists := event[g.target]; !exists {
					event[g.target] = target
				}
			}

		}
		return event, nil
	}
	return event, errors.New("grokerFilter match all not patterns")
}
