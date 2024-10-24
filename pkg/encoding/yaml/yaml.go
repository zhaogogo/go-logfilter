package yaml

import (
	"gopkg.in/yaml.v3"
)

// Name is the name registered for the yaml Code.
const Name = "yaml"

// Code is a Codec implementation with yaml.
type Code struct{}

func (Code) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (Code) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (Code) Name() string {
	return Name
}
