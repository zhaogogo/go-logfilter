package file

import (
	"dario.cat/mergo"
	"testing"
)

func TestWatch(t *testing.T) {
	source := NewSource("/Users/zhaoqiang/Documents/owner项目/go-logfilter/config")
	kvs, err := source.Load()
	if err != nil {
		t.Error(err)
		return
	}
	for _, kv := range kvs {
		t.Log("load", kv.Key, kv.Format, string(kv.Value))
	}

	watch, err := source.Watch()
	if err != nil {
		t.Error(err)
		return
	}
	for {
		keyValues, err := watch.Next()
		if err != nil {
			t.Error(err)
			return
		}
		for _, kv := range keyValues {
			t.Log("watch", kv.Key, kv.Format, string(kv.Value))
		}
	}
}

func TestMerger(t *testing.T) {
	src := map[string]interface{}{"a": []interface{}{"1", "2"}}
	dest := map[string]interface{}{"a": []interface{}{"b", "2"}}
	err := mergo.Merge(&dest, src)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v\n", dest)
}
