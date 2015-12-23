package test

import (
	"fmt"

	"github.com/ningjh/memcached"
	"github.com/ningjh/memcached/config"
)

var Mcpool *memcached.MemcachedClient4B

func init() {
	var conf = config.New()

	conf.Servers = []string{"10.210.215.245:11211"}
	conf.ReadTimeout = 1000
	conf.WriteTimeout = 1000
	conf.InitConns = 50
	conf.NumberOfReplicas = 25

	var err error
	Mcpool, err = memcached.NewMemcachedClient4B(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func TestMcpool() {
	// 保存数据
	//var element = &common.Element{
	//	Key:     "test-golang",
	//	Flags:   0,
	//	Exptime: 3600,
	//	Value:   []byte("test mc pool"),
	//}
	//err = memcachedClient.Add(element)

	// 单个key获取数据
	item, err := Mcpool.Get("test-golang")
	if err != nil {
		fmt.Println(err)
		return
	}
	key := item.Key()
	value := item.Value()
	flags := item.Flags()
	cas := item.Cas()

	fmt.Println(key, string(value), flags, cas)

	// 多个key获取数据
	//keys := []string{"abc", "def", "ghi"}
	//items, err := memcachedClient.GetArray(keys)
	//if err == nil {
	//	for _, key := range keys {
	//		if item, ok := items[key]; ok {
	//			key := item.Key()
	//			value := item.Value()
	//			flags := item.Flags()
	//			cas := item.Cas()
	//		}
	//	}
	//}
}
