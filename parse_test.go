package config

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	Config := New("./config.ini",true).SetDefaultGroup("mysql")
	if !Config.Load() {
		t.Error(Config.LastErr())
	}



	var maps = map[string]interface{} {
		"id":"1",
		"memcached.id": "4",
		"redis.id" : "3",
		"redis.conv": "5",
		"memcached.conv":nil,
		"memcached.arr":[]string{"1","3"},
	}

	for key,expected := range maps {
		res := Config.Get(key)

		if !reflect.DeepEqual(res,expected) {
			t.Errorf("结果值为:%v 期望值为%v",res,expected)
		}
	}
}

func TestGetGroup(t *testing.T) {
	Config := New("./config.ini",true).SetDefaultGroup("prod")
	if !Config.Load() {
		t.Error(Config.LastErr())
	}



	var maps = map[string]interface{} {
		"mysql": map[string]interface{} {"id":"1",},
		"memcached": map[string]interface{} {"id":"4","arr":[]string{"1","3"}},
		"redis":map[string]interface{} {"id":"3","conv":"5"},
	}

	for key,expected := range maps {
		res,ok := Config.GetGroup(key)
		if !ok {
			t.Error(Config.LastErr())
		}
		if !reflect.DeepEqual(res,expected) {
			t.Errorf("结果值为:%v 期望值为%v",res,expected)
		}
	}
}
