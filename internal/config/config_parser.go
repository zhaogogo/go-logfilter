package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zhaogogo/go-logfilter/pkg/encoding"
	"gopkg.in/yaml.v3"
)

type Config map[string]interface{}

type Parser interface {
	parse(filepath string) (map[string]interface{}, error)
}

func ParseConfig(source Source) (map[string][]interface{}, error) {
	kvs, err := source.Load()
	if err != nil {
		return nil, err
	}
	conf := make(map[string][]interface{})
	for _, kv := range kvs {
		c := make(map[string]interface{})
		code, ok := encoding.GetCodec(kv.Format)
		if !ok {
			log.Info().Msgf("未注册的文件类型 key=%s format=%s", kv.Key, kv.Format)
			continue
		}
		log.Info().Msgf("解析文件 %s", kv.Key)
		if err := code.Unmarshal(kv.Value, c); err != nil {
			return nil, errors.Wrapf(err, "解析文件%s失败", kv.Key)
		}
		conf = MergeConfig(conf, c)
	}
	confy, _ := yaml.Marshal(conf)
	fmt.Printf("合并后配置文件为: %s\n", confy)
	return conf, nil
}

func MergeConfig(dest map[string][]interface{}, src map[string]interface{}) map[string][]interface{} {
	if inputs, ok := src["inputs"]; ok {
		if v, ok := inputs.([]interface{}); ok {
			dest["inputs"] = append(dest["inputs"], v...)
		}
	}
	if filters, ok := src["filters"]; ok {
		if v, ok := filters.([]interface{}); ok {
			dest["filters"] = append(dest["filters"], v...)
		}
	}
	if outputs, ok := src["outputs"]; ok {
		if v, ok := outputs.([]interface{}); ok {
			dest["outputs"] = append(dest["outputs"], v...)
		}
	}
	return dest
}
