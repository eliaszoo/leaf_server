package data

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"leaf_server/conf"

	"github.com/name5566/leaf/log"
)

const (
	TABLENAME_BASIC = "test.json"
)

var CFG_PATH string

type TABLE_ITEM interface {
	GetID() interface{}
}

type TABLE map[interface{}]TABLE_ITEM

var DataMap map[string]TABLE

func init() {
	CFG_PATH = conf.Server.ConfPath
	DataMap = make(map[string]TABLE)

	readcfg(TABLENAME_BASIC, &[]Test{})
}

func GetTable(name string) TABLE {
	return DataMap[name]
}

func GetTableItem(name string, key interface{}) TABLE_ITEM {
	table, ok := DataMap[name]
	if !ok {
		return nil
	}

	return table[key]
}

func readcfg(name string, array interface{}) error {
	arrayv := reflect.ValueOf(array)
	if arrayv.Kind() != reflect.Ptr || arrayv.Elem().Kind() != reflect.Slice {
		log.Error("readfile: output isn't ptr and kind isn't slice ", name)
		return nil
	}

	data, err := ioutil.ReadFile(CFG_PATH + name)
	if nil != err {
		log.Error("read file failed:", name, " err:", err.Error())
		return err
	}

	err = json.Unmarshal(data, array)
	if nil != err {
		log.Error("json unmarshal failed:", name, " err:", err.Error())
		return err
	}

	table := make(map[interface{}]TABLE_ITEM)
	arraye := arrayv.Elem()
	for i := 0; i < arraye.Len(); i++ {
		v := arraye.Index(i).Interface()
		item := v.(TABLE_ITEM)
		table[item.GetID()] = item
	}
	DataMap[name] = table
	return nil
}
