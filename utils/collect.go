package utils

import (
	"fmt"
	"github.com/chenhg5/collection"
	"github.com/livegoplayer/go_helper/structs_map"
	"reflect"
)

type Collect = collection.Collection

type Collection struct {
	Collect
}

type dataList = []map[string]interface{}

type DefaultFunc = func(map[string]interface{}) interface{}

func NewCollect(st interface{}, tagName string) Collection {
	var v interface{}
	v = st
	mapArray, ok := v.(dataList)
	if !ok {
		slice := reflect.ValueOf(v)
		if slice.Kind() == reflect.Ptr {
			slice = slice.Elem()
		}
		for i := 0; i < slice.Len(); i++ {
			mapArray = append(mapArray, structs_map.NewStructMap(slice.Index(i).Addr().Elem().Interface(), tagName).Map())
		}
	}
	return Collection{collection.Collect(mapArray)}
}

func (c Collection) WithGroupBy(asKey string, list dataList, localKey, localKeyTagName string, argus ...interface{}) Collection {
	foreign, defVal := defaultArgus(argus, localKey, dataList{})
	srcData := c.ToMapArray()
	withData := NewCollect(list, localKeyTagName).GroupBy(foreign).ToMap()
	for _, item := range srcData {
		v, ok := DeepGet(item, localKey)
		if ok {
			tv, ok := withData[fmt.Sprintf("%v", v)]
			if ok {
				item[asKey] = tv
				continue
			}
		}
		cb, ok := defVal.(DefaultFunc)
		if ok {
			item[asKey] = cb(item)
		} else {
			item[asKey] = defVal
		}
	}
	return NewCollect(srcData, "")
}

func (c Collection) WithKeyBy(asKey string, list dataList, localKey, localKeyTagName string, argus ...interface{}) Collection {
	foreign, defVal := defaultArgus(argus, localKey, nil)
	srcData := c.ToMapArray()
	withData := NewCollect(list, localKeyTagName).KeyBy(foreign).ToMap()
	for _, item := range srcData {
		v, ok := DeepGet(item, localKey)
		if ok {
			tv, ok := withData[fmt.Sprintf("%v", v)]
			if ok {
				item[asKey] = tv.(dataList)[0]
				continue
			}
		}
		cb, ok := defVal.(DefaultFunc)
		if ok {
			item[asKey] = cb(item)
		} else {
			item[asKey] = defVal
		}
	}
	return NewCollect(srcData, "")
}

func (c Collection) WithPut(asKey string, fromKey string, list dataList, localKey, localKeyTagName string, argus ...interface{}) Collection {
	foreign, defVal := defaultArgus(argus, localKey, "")
	srcData := c.ToMapArray()
	withData := NewCollect(list, localKeyTagName).KeyBy(foreign).ToMap()
	for _, item := range srcData {
		v, ok := DeepGet(item, localKey)
		if ok {
			tv, ok := withData[fmt.Sprintf("%v", v)]
			if ok {
				if v, ok := tv.(dataList)[0][fromKey]; ok {
					DeepMustSet(item, asKey, v)
					continue
				}
			}
		}
		cb, ok := defVal.(DefaultFunc)
		if ok {
			DeepMustSet(item, asKey, cb(item))
		} else {
			DeepMustSet(item, asKey, defVal)
		}
	}
	return NewCollect(srcData, "")
}

func (c Collection) WithConst(asKey string, val interface{}) Collection {
	srcData := c.ToMapArray()
	for _, item := range srcData {
		item[asKey] = val
	}
	return NewCollect(srcData, "")
}

func (c Collection) EachPut(asKey string, cb func(map[string]interface{}) interface{}) Collection {
	srcData := c.ToMapArray()
	for _, item := range srcData {
		item[asKey] = cb(item)
	}
	return NewCollect(srcData, "")
}

func (c Collection) One(cb func(idx int, t map[string]interface{}) bool) Collection {
	srcData := c.ToMapArray()
	for i, item := range srcData {
		if cb(i, item) {
			return NewCollect([]H{item}, "")
		}
	}
	return NewCollect(struct{}{}, "")
}

func (c Collection) Each(cb func(idx int, t map[string]interface{})) Collection {
	srcData := c.ToMapArray()
	for i, item := range srcData {
		cb(i, item)
	}
	return NewCollect(srcData, "")
}

func defaultArgus(argus []interface{}, foreign string, val interface{}) (string, interface{}) {
	switch len(argus) {
	case 0:
		return foreign, val
	case 1:
		return argus[0].(string), val
	case 2:
		return argus[0].(string), argus[1]
	}
	panic("参数错误")
}
