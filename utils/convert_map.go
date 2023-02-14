package utils

import (
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

// Merge AppendStruct 将give的结构体中，同名的成员赋值给
func Merge(received Any, give Any) {
	givenType := reflect.TypeOf(give).Elem()
	givenVal := reflect.ValueOf(give).Elem()

	receivedType := reflect.TypeOf(received).Elem()
	receivedVal := reflect.ValueOf(received).Elem()

	for i := 0; i < givenType.NumField(); i++ {
		givenName := givenType.Field(i).Name
		for j := 0; j < receivedType.NumField(); j++ {
			if givenName == receivedType.Field(j).Name {
				receivedVal.Field(j).Set(givenVal.Field(i))
				break
			}
		}
	}
}

// Struct2Struct 必须dest传入指针
func Struct2Struct(src, dest interface{}) {
	ToStruct(structs.Map(src), dest)
}

func anyToStruct(m interface{}, s interface{}, tagNames ...string) {
	tagName := "json"
	if len(tagNames) > 0 {
		tagName = tagNames[0]
	}

	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           s,
		TagName:          tagName,
		WeaklyTypedInput: true,
	}
	decoder, _ := mapstructure.NewDecoder(config)
	err := decoder.Decode(m)
	if err != nil {
		panic(err)
	}
}

func ToStruct(m interface{}, s interface{}, tagNames ...string) {
	if str, ok := m.(string); ok {
		jsonToStruct(str, s, tagNames...)
		return
	}
	// 非字符串内容，由他自己处理
	anyToStruct(m, s, tagNames...)
}

// jsonToStruct 将json字符串，复制给s结构体中同名json-tag的字段，s必须为结构体指针
func jsonToStruct(json string, s interface{}, tagNames ...string) {
	m := JsonDecodeSafe(json)
	if m == nil {
		return
	}
	ToStruct(m, s, tagNames...)
}

func JsonDecodeSafe(jsonStr string) map[string]interface{} {
	var mapResult map[string]interface{}
	ds := json.NewDecoder(strings.NewReader(jsonStr))
	ds.UseNumber()
	_ = ds.Decode(&mapResult)
	return mapResult
}

func JsonDecodes(jsonStr string) (map[string]interface{}, []map[string]interface{}) {
	if jsonStr[0] == '[' {
		mapSlice := make([]map[string]interface{}, 0)
		_ = json.Unmarshal([]byte(jsonStr), &mapSlice)
		return nil, mapSlice
	}
	var mapResult map[string]interface{}
	_ = json.Unmarshal([]byte(jsonStr), &mapResult)
	return mapResult, nil
}

func JsonDecodeAny(jsonStr string) interface{} {
	var a Any
	_ = json.Unmarshal([]byte(jsonStr), &a)
	return a
}

func MapKeys(item map[string]interface{}) []string {
	ks := make([]string, 0)
	for k := range item {
		ks = append(ks, k)
	}
	return ks
}

// ToMap converts a struct to a map using the struct's tags.
//
// ToMap uses tags on struct fields to decide which fields to add to the
// returned map.
func ToMap(in interface{}) map[string]interface{} {
	return NewCollect(in, "json").ToMap()
}

func ToMapArray(in interface{}) []map[string]interface{} {
	return NewCollect(in, "json").ToMapArray()
}

func DeepGetMust(m map[string]interface{}, key string) interface{} {
	v, ok := DeepGet(m, key)
	if !ok {
		panic("没有该值")
	}
	return v
}

func DeepGetShould(m map[string]interface{}, key string) interface{} {
	v, _ := DeepGet(m, key)
	return v
}

func DeepMustSet(m map[string]interface{}, key string, val interface{}) {
	succ := DeepSet(m, key, val)
	if !succ {
		panic("设置失败")
	}
}

func DeepGet(m map[string]interface{}, key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	var ret interface{} = m
	for _, k := range keys {
		v, ok := ret.(map[string]interface{})
		if !ok {
			return nil, false
		}
		ret, ok = v[k]
		if !ok {
			return nil, false
		}
	}
	return ret, true
}

func DeepSet(m map[string]interface{}, key string, val interface{}) bool {
	keys := strings.Split(key, ".")
	var ret = m
	for i, k := range keys {
		if i == len(keys)-1 {
			ret[k] = val
			return true
		}
		tmpRet, ok := ret[k]
		if !ok {
			return false
		}
		tmpRetm, ok := tmpRet.(map[string]interface{})
		if !ok {
			return false
		}
		ret = tmpRetm
	}
	return false
}
