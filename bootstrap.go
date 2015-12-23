package main

import (
	"weibo.com/hotweibo/lib/config"
	"weibo.com/hotweibo/service"
)

func init() {
	config.LoadDir("config")
	basic.InitMCPool()
	basic.InitHotMCPool()
	basic.InitBasicDbPool()
	basic.InitDbPool()
	basic.InitRedisPool()
	basic.InitRedis7474Pool()
	basic.InitRedis7475Pool()
	basic.InitHttpPool()
}
