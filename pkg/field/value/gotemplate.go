package valuerender

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

var GOHANGOUT_TYPE_UNKNOWN_ERROR error = errors.New("field type unknown, it must be of json.Number|Int64|Int|int8")

var ErrNotFloat64 error = errors.New("Only float64 type value could be calculated")
var ErrNotInt64 error = errors.New("Only int64 type value could be calculated")

var funcMap = template.FuncMap{}

type TemplateValueRender struct {
	tmpl *template.Template
}

func NewTemplateValueRender(t string) *TemplateValueRender {
	tmpl, err := template.New(t).Funcs(funcMap).Parse(t)
	if err != nil {
		log.Fatal().Msgf("could not parse template %s:%s", t, err)
	}
	return &TemplateValueRender{
		tmpl: tmpl,
	}
}

func init() {
	for k, v := range sprig.FuncMap() {
		funcMap[k] = v
	}

	funcMap["compare"] = strings.Compare
	funcMap["contains"] = strings.Contains
	funcMap["containsAny"] = strings.ContainsAny
	funcMap["hasprefix"] = strings.HasPrefix
	funcMap["hassuffix"] = strings.HasSuffix
	funcMap["replace"] = strings.Replace
	funcMap["timeFormat"] = func(t interface{}, srcFormat string, format string) (string, error) {
		switch t.(type) {
		case time.Time:
			tt := t.(time.Time)
			return tt.Format(format), nil
		case string:
			t, err := time.Parse(srcFormat, t.(string))
			if err != nil {
				log.Err(err).Msgf("parse time failed")
				return "", err
			}
			return t.Format(format), nil
		default:
			return "", errors.New("go template 不支持的输入类型")
		}
	}

	funcMap["now"] = func() int64 { return time.Now().UnixNano() / 1000000 }
	funcMap["timestamp"] = func(event map[string]interface{}) int64 {
		timestamp := event["@timestamp"]
		if timestamp == nil {
			return 0
		}
		if reflect.TypeOf(timestamp).String() == "time.Time" {
			return timestamp.(time.Time).UnixNano() / 1000000
		}
		return 0
	}

	funcMap["before"] = func(event map[string]interface{}, s string) bool {
		timestamp := event["@timestamp"]
		if timestamp == nil || reflect.TypeOf(timestamp).String() != "time.Time" {
			return false
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			log.Err(err)
			return false
		}
		dst := time.Now().Add(d)
		return timestamp.(time.Time).Before(dst)
	}

	funcMap["after"] = func(event map[string]interface{}, s string) bool {
		timestamp := event["@timestamp"]
		if timestamp == nil || reflect.TypeOf(timestamp).String() != "time.Time" {
			return false
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			log.Err(err)
			return false
		}
		dst := time.Now().Add(d)
		return timestamp.(time.Time).After(dst)
	}

	funcMap["plus"] = func(x, y interface{}) (float64, error) {
		if xf, ok := x.(float64); ok {
			if yf, ok := y.(float64); ok {
				return xf + yf, nil
			}
		}
		return 0, ErrNotFloat64
	}

	funcMap["minus"] = func(x, y interface{}) (float64, error) {
		if xf, ok := x.(float64); ok {
			if yf, ok := y.(float64); ok {
				return xf - yf, nil
			}
		}
		return 0, ErrNotFloat64
	}
	funcMap["multiply"] = func(x, y interface{}) (float64, error) {
		if xf, ok := x.(float64); ok {
			if yf, ok := y.(float64); ok {
				return xf * yf, nil
			}
		}
		return 0, ErrNotFloat64
	}
	funcMap["divide"] = func(x, y interface{}) (float64, error) {
		if xf, ok := x.(float64); ok {
			if yf, ok := y.(float64); ok {
				return xf / yf, nil
			}
		}
		return 0, ErrNotFloat64
	}
	funcMap["mod"] = func(x, y interface{}) (int64, error) {
		if xf, ok := x.(int64); ok {
			if yf, ok := y.(int64); ok {
				return xf % yf, nil
			}
		}
		return 0, ErrNotInt64
	}
}

// always return string
func (r *TemplateValueRender) Render(event map[string]interface{}) interface{} {
	b := bytes.NewBuffer(nil)
	if err := r.tmpl.Execute(b, event); err != nil {
		return nil
	}
	res := b.String()
	if res == "" {
		return nil
	}
	return res
}

func convertToInt(x interface{}) (int, error) {
	if reflect.TypeOf(x).String() == "json.Number" {
		b, _ := x.(json.Number).Int64()
		return int(b), nil
	} else if reflect.TypeOf(x).Kind() == reflect.Int64 {
		return int(x.(int64)), nil
	} else if reflect.TypeOf(x).Kind() == reflect.Int {
		return x.(int), nil
	} else if reflect.TypeOf(x).Kind() == reflect.Int8 {
		return int(x.(int8)), nil
	}
	return 0, GOHANGOUT_TYPE_UNKNOWN_ERROR
}
