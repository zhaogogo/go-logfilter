package condition

import (
	"math/rand"
	"time"
)

type RandomCondition struct {
	value int
}

func NewRandomCondition(value int) *RandomCondition {
	rand.Seed(time.Now().UnixNano())
	return &RandomCondition{value}
}

func (c *RandomCondition) Pass(event map[string]interface{}) bool {
	return rand.Intn(c.value) == 0
}
