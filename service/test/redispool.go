package test

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	pool          *redis.Pool
	redisServer   = "10.210.215.245:6379"
	redisPassword = ""
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newPool(redisServer, redisPassword)
}

func TestRedis() {
	conn := pool.Get()
	defer conn.Close()
	n, err := conn.Do("get", "test-redigo")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(n.([]byte)))
	}
}
