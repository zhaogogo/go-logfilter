package condition

import (
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
)

type BeforeCondition struct {
	d time.Duration
}

func NewBeforeCondition(value string) *BeforeCondition {
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatal().Msgf("could not parse %s to duration: %s", value, err)
	}
	return &BeforeCondition{d}
}

func (c *BeforeCondition) Pass(event map[string]interface{}) bool {
	timestamp := event["@timestamp"]
	if timestamp == nil || reflect.TypeOf(timestamp).String() != "time.Time" {
		return false
	}
	return timestamp.(time.Time).Before(time.Now().Add(c.d))
}

type AfterCondition struct {
	d time.Duration
}

func NewAfterCondition(value string) *AfterCondition {
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatal().Msgf("could not parse %s to duration: %s", value, err)
	}
	return &AfterCondition{d}
}

func (c *AfterCondition) Pass(event map[string]interface{}) bool {
	timestamp := event["@timestamp"]
	if timestamp == nil || reflect.TypeOf(timestamp).String() != "time.Time" {
		return false
	}
	return timestamp.(time.Time).After(time.Now().Add(c.d))
}
