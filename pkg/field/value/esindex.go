package valuerender

// used for ES indexname template

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type field struct {
	literal  bool
	date     bool
	raw      string
	value    string // used in datetime %{+} and literal
	location *time.Location

	mv *MultiLevelValueRender
}

func (f *field) render(event map[string]interface{}, location *time.Location) (interface{}, error) {
	if f.literal {
		return f.value, nil
	}

	if f.date {
		if t, ok := event["@timestamp"]; ok {
			return dateFormat(t, f.value, location)
		} else {
			return dateFormat(time.Now(), f.value, location)
		}
	}
	v := f.mv.Render(event)
	return v, nil
}

type IndexRender struct {
	fields   []*field
	location *time.Location
}

// getAllFields("%{@metadata}{kafka}{topic}") => ["@metadata","kafka","topic"]
func getAllFields(s string) []string {
	fields := make([]string, 0)
	r, _ := regexp.Compile(`{(.*?)}`)
	for _, v := range r.FindAll([]byte(s), -1) {
		fields = append(fields, string(v[1:len(v)-1]))
	}
	return fields
}

func NewIndexRender(t string) *IndexRender {
	r, _ := regexp.Compile(`%({.*?})+`) //%{+YYYY.MM.dd}
	fields := make([]*field, 0)
	lastPos := 0
	for _, loc := range r.FindAllStringIndex(t, -1) {
		s, e := loc[0], loc[1]
		if lastPos != s {
			fields = append(fields, &field{
				literal: true,
				raw:     t[lastPos:s],
				value:   t[lastPos:s],
			})
		}

		if t[s+2] == '+' {
			fields = append(fields, &field{
				literal: false,
				date:    true,
				raw:     t[s:e],
				value:   t[s+3 : e-1],
			})
		} else {
			fields = append(fields, &field{
				literal: false,
				date:    false,
				raw:     t[s:e],
				mv:      NewMultiLevelValueRender(getAllFields(t[s+1 : e])),
			})
		}

		lastPos = e
	}

	if lastPos < len(t) {
		fields = append(fields, &field{
			literal: true,
			date:    false,
			raw:     t[lastPos:],
			value:   t[lastPos:],
		})
	}
	return &IndexRender{fields, time.UTC}
}

// SetTimeLocation parse `location` to time.Location ans set it as its member.
// use this location to format time string
func (r *IndexRender) SetTimeLocation(loc string) {
	location, err := time.LoadLocation(loc)
	if err != nil {
		log.Fatal().Msgf("invalid localtion: %s", loc)
	}
	r.location = location
}

func dateFormat(t interface{}, format string, location *time.Location) (string, error) {
	if t1, ok := t.(time.Time); ok {
		return t1.In(location).Format(format), nil
	}
	if reflect.TypeOf(t).String() == "json.Number" {
		t1, err := t.(json.Number).Int64()
		if err != nil {
			return format, err
		}
		return time.Unix(t1/1000, t1%1000*1000000).In(location).Format(format), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.Int {
		t1 := int64(t.(int))
		return time.Unix(t1/1000, t1%1000*1000000).In(location).Format(format), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.Int64 {
		t1 := t.(int64)
		return time.Unix(t1/1000, t1%1000*1000000).In(location).Format(format), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.String {
		t1, e := time.Parse(time.RFC3339, t.(string))
		if e != nil {
			return format, e
		}
		return t1.In(location).Format(format), nil
	}
	return format, errors.New("could not tell the type timestamp field belongs to")
}

func (r *IndexRender) Render(event map[string]interface{}) interface{} {
	fields := make([]string, len(r.fields))
	for i, f := range r.fields {
		if v, err := f.render(event, r.location); err != nil {
			log.Err(err).Msgf("esIndex render failed")
			fields[i] = f.value
		} else {
			res, ok := v.(string)
			if !ok {
				log.Error().Msgf("esIndex render asset failed, got: %T  %#v", v, v)
				fields[i] = f.raw
			} else {
				fields[i] = res
			}
		}
	}
	return strings.Join(fields, "")
}
