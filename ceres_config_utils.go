package CeresConfigEtcd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"gopkg.in/yaml.v3"
	"strings"
)

var (
	Marshals = map[string]Marshal{
		"json": json.Marshal,
		"xml":  xml.Marshal,
		"yaml": yaml.Marshal,
		"yml":  yaml.Marshal,
		"toml": func(v interface{}) ([]byte, error) {
			b := bytes.NewBuffer(nil)
			defer b.Reset()
			err := toml.NewEncoder(b).Encode(v)
			if err != nil {
				return nil, err
			}
			return b.Bytes(), nil
		},
	}
	Unmarshals = map[string]Unmarshal{
		"json": json.Unmarshal,
		"xml":  xml.Unmarshal,
		"yaml": yaml.Unmarshal,
		"yml":  yaml.Unmarshal,
		"toml": toml.Unmarshal,
	}
)

type Marshal func(v interface{}) ([]byte, error)
type Unmarshal func(data []byte, v interface{}) error

// makeMapData 把[]byte根据配置转为map
func makeMapData(kv []*mvccpb.KeyValue, fn Unmarshal, trimPrefix string) map[string]interface{} {
	data := make(map[string]interface{})
	for _, value := range kv {
		data = modifyMapData("put", trimPrefix, data, value, fn)
	}
	return data
}

// modifyMapData 调整数据
func modifyMapData(op, trimPrefix string, data map[string]interface{}, kv *mvccpb.KeyValue, fn Unmarshal) map[string]interface{} {
	// 删除前缀，例如：/ceres/config/etcd/default,操作后的为：etcd/default
	key := strings.TrimPrefix(strings.TrimPrefix(string(kv.Key), trimPrefix), "/")
	// 分割为["etcd","default"]
	keys := strings.Split(key, "/")
	// 序列化数据
	var value interface{}
	_ = fn(kv.Value, &value)

	if len(keys) > 0 && len(keys) == 1 {
		switch op {
		case "delete":
			data = make(map[string]interface{})
		default:
			v, ok := value.(map[string]interface{})
			if ok {
				data = v
			}
		}
		return data
	}

	tempData := data
	for i, k := range keys {
		// 先判断该key是否已经存在值
		kData, ok := data[k].(map[string]interface{})
		if !ok {
			kData = make(map[string]interface{})
			tempData[k] = kData
		}
		// 如果是最后一个key，则设置数据
		if len(keys)-1 == i {
			switch op {
			case "delete":
				delete(tempData, k)
			default:
				tempData[k] = value
			}
		} else {
			tempData = kData
		}
	}
	return data
}
