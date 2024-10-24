package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/pkg/field"
	"math"
	"reflect"
	"strconv"
	"time"
)

type DateParser interface {
	Parse(interface{}) (time.Time, error)
}

type FormatParser struct {
	format   string
	location *time.Location
	addYear  bool
}

var MustStringTypeError = errors.New("timestamp field must be string")

func (dp *FormatParser) Parse(t interface{}) (time.Time, error) {
	var (
		rst time.Time
		err error
	)
	fmt.Println("FormatParser -->", t)
	if reflect.TypeOf(t).String() != "string" {
		return rst, MustStringTypeError
	}
	var value string
	if dp.addYear {
		value = fmt.Sprintf("%d%s", time.Now().Year(), t.(string))
	} else {
		value = t.(string)
	}
	if dp.location == nil {
		return time.Parse(dp.format, value)
	}
	rst, err = time.ParseInLocation(dp.format, value, dp.location)
	if err != nil {
		return rst, err
	}
	return rst.UTC(), nil
}

type UnixParser struct{}

func (p *UnixParser) Parse(t interface{}) (time.Time, error) {
	var (
		rst time.Time
	)
	if v, ok := t.(json.Number); ok {
		t1, err := v.Int64()
		if err != nil {
			return rst, err
		}
		return time.Unix(t1, 0), nil
	}

	if v, ok := t.(string); ok {
		t1, err := strconv.Atoi(v)
		if err != nil {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return rst, err
			}
			t1 := math.Floor(f)
			return time.Unix(int64(t1), int64(1000000000*(f-t1))), nil
		}
		return time.Unix(int64(t1), 0), nil
	}

	if t1, ok := t.(int); ok {
		return time.Unix(int64(t1), 0), nil
	}
	if t1, ok := t.(int64); ok {
		return time.Unix(t1, 0), nil
	}
	return rst, fmt.Errorf("%s unknown type:%s", t, reflect.TypeOf(t).String())
}

type UnixMSParser struct{}

func (p *UnixMSParser) Parse(t interface{}) (time.Time, error) {
	var (
		rst time.Time
	)
	if reflect.TypeOf(t).String() == "json.Number" {
		t1, err := t.(json.Number).Int64()
		if err != nil {
			return rst, err
		}
		return time.Unix(t1/1000, t1%1000*1000000), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.String {
		t1, err := strconv.Atoi(t.(string))
		if err != nil {
			return rst, err
		}
		t2 := int64(t1)
		return time.Unix(t2/1000, t2%1000*1000000), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.Int {
		t1 := int64(t.(int))
		return time.Unix(t1/1000, t1%1000*1000000), nil
	}
	if reflect.TypeOf(t).Kind() == reflect.Int64 {
		t1 := t.(int64)
		return time.Unix(t1/1000, t1%1000*1000000), nil
	}
	return rst, fmt.Errorf("%s unknown type:%s", t, reflect.TypeOf(t).String())
}

func getDateParser(format string, l *time.Location, addYear bool) DateParser {
	if format == "UNIX" {
		return &UnixParser{}
	}
	if format == "UNIX_MS" {
		return &UnixMSParser{}
	}
	if format == "RFC3339" {
		return &FormatParser{time.RFC3339, l, addYear}
	}
	return &FormatParser{format, l, addYear}
}

type DateFilter struct {
	config       map[interface{}]interface{}
	dateParsers  []DateParser
	overwrite    bool
	src          string
	srcVR        field.ValueRender
	target       string
	targetFS     field.FieldSetter
	targetFormat string
}

func New(config map[interface{}]interface{}) *DateFilter {
	plugin := &DateFilter{
		config:      config,
		overwrite:   true,
		dateParsers: make([]DateParser, 0),
	}

	if overwrite, ok := config["overwrite"]; ok {
		plugin.overwrite = overwrite.(bool)
	}

	if srcValue, ok := config["src"]; ok {
		plugin.src = srcValue.(string)
	} else {
		log.Fatal().Msg("src must be set in date filter plugin")
	}
	plugin.srcVR = field.GetValueRender2(plugin.src)

	if targetI, ok := config["target"]; ok {
		plugin.target = targetI.(string)
	} else {
		plugin.target = "@timestamp"
	}
	plugin.targetFS = field.NewFieldSetter(plugin.target)
	if targetFormat, ok := config["targetFormat"]; ok {
		plugin.targetFormat = targetFormat.(string)
	}
	var (
		location *time.Location
		addYear  bool = false
		err      error
	)
	if locationI, ok := config["location"]; ok {
		location, err = time.LoadLocation(locationI.(string))
		if err != nil {
			log.Fatal().Msgf("load location error:%s", err)
		}
	} else {
		location = nil
	}
	if addYearI, ok := config["add_year"]; ok {
		addYear = addYearI.(bool)
	}
	if formats, ok := config["formats"]; ok {
		for _, formatI := range formats.([]interface{}) {
			plugin.dateParsers = append(plugin.dateParsers, getDateParser(formatI.(string), location, addYear))
		}
	} else {
		log.Fatal().Msg("formats must be set in date filter plugin")
	}

	return plugin
}

func (plugin *DateFilter) Filter(event map[string]interface{}) (map[string]interface{}, bool) {
	inputI := plugin.srcVR.Render(event)
	if inputI == nil {
		return event, false
	}
	fmt.Println("--->", inputI)
	for _, dp := range plugin.dateParsers {
		t, err := dp.Parse(inputI)
		if err == nil {
			fmt.Println("---->", plugin.targetFormat)
			if plugin.targetFormat != "" {
				timeText := t.Format(plugin.targetFormat)
				fmt.Printf("解析时间：%s\n", timeText)
				event = plugin.targetFS.SetField(event, timeText, plugin.overwrite)
				return event, true
			}
			event = plugin.targetFS.SetField(event, t, plugin.overwrite)
			return event, true
		}
	}
	return event, false
}
