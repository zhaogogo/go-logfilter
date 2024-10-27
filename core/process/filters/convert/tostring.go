package convert

import "encoding/json"

type StringConverter struct{}

func (c *StringConverter) convert(v interface{}) (interface{}, error) {
	if r, ok := v.(json.Number); ok {
		return r.String(), nil
	}

	if r, ok := v.(string); ok {
		return r, nil
	}

	jsonString, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return string(jsonString), nil
}
