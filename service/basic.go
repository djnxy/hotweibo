package basic

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ningjh/memcached"
	"github.com/ningjh/memcached/config"
)

var (
	Mcpool        *memcached.MemcachedClient4B
	HotMcpool     *memcached.MemcachedClient4B
	Dbpool        *sql.DB
	BasicDbpool   *sql.DB
	Redispool     *redis.Pool
	Redis7474pool *redis.Pool
	Redis7475pool *redis.Pool
	Httppool      *http.Client
)

func InitRedis7474Pool() {
	server := "host:port"
	password := ""
	Redis7474pool = &redis.Pool{
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

func InitRedis7475Pool() {
	server := "host:port"
	password := ""
	Redis7475pool = &redis.Pool{
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

func InitRedisPool() {
	server := "host:port"
	password := ""
	Redispool = &redis.Pool{
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

func InitDbPool() {
	var err error
	Dbpool, err = sql.Open("mysql", "root:123456@tcp(host:port)/test?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	Dbpool.SetMaxOpenConns(0)
	Dbpool.SetMaxIdleConns(50)
	Dbpool.Ping()
}

func InitBasicDbPool() {
	var err error
	BasicDbpool, err = sql.Open("mysql", "root:test@tcp(host:port)/hotmblog?charset=utf8")
	if err != nil {
		fmt.Println(err)
		return
	}
	BasicDbpool.SetMaxOpenConns(0)
	BasicDbpool.SetMaxIdleConns(50)
	BasicDbpool.Ping()
}

func InitHotMCPool() {
	var conf = config.New()

	conf.Servers = []string{"host:port"}
	conf.ReadTimeout = 1000
	conf.WriteTimeout = 1000
	conf.InitConns = 50
	conf.NumberOfReplicas = 25

	var err error
	HotMcpool, err = memcached.NewMemcachedClient4B(conf)
	if err != nil {
		fmt.Println(err)
	}
}

func InitMCPool() {
	var conf = config.New()

	conf.Servers = []string{"host:port"}
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

func InitHttpPool() {
	Httppool = &http.Client{

		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*1)
				if err != nil {
					fmt.Println("dail timeout", err)
					return nil, err
				}
				return c, nil

			},
			MaxIdleConnsPerHost:   50,
			ResponseHeaderTimeout: time.Second * 1,
		},
	}
}
