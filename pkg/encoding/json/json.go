package json

import (
	"encoding/json"
)

// Name is the name registered for the yaml Code.
const Name = "json"

// Code is a Codec implementation with yaml.
type Code struct{}

func (Code) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (Code) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (Code) Name() string {
	return Name
}
