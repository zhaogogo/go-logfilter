package config

type Config map[string]interface{}

type Parser interface {
	parse(filepath string) (map[string]interface{}, error)
}

func ParseConfig(source Source) (map[string]interface{}, error) {
	kvs, err := source.Load()
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		kv.Format
	}
	return nil, nil
}
