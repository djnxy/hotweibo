package test

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

func TestMc() {
	mc := memcache.New("10.210.215.245:11211")

	value, err := mc.Get("test-golang")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(value.Key, string(value.Value))
	}

}
