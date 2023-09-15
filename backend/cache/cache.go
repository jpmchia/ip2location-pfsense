package cache

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson/v4"
)

type RedisCache struct {
	RedisHostPort string
	RedisDb       int
	RedisAuth     string
	RedisPass     string
}

type Instance struct {
	RedisInstance RedisCache
	RedisPool     *redis.Pool
	Rh            *rejson.Handler
	ResultSet     int
}

func init() {

}

func NewPool(ri RedisCache) *redis.Pool {

	log.Printf("Creating Redis client pool for database %d on %s.\n", ri.RedisDb, ri.RedisHostPort)

	return &redis.Pool{

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ri.RedisHostPort, redis.DialDatabase(ri.RedisDb), redis.DialPassword(ri.RedisPass))
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("SELECT", ri.RedisHostPort); err != nil {
				err = c.Close()
				return nil, err
			}
			if _, err := c.Do("AUTH", ri.RedisPass); err != nil {
				err = c.Close()
				return nil, err
			}
			return c, err
		},
		DialContext: nil,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         5,
		MaxActive:       10,
		IdleTimeout:     240 * time.Second,
		Wait:            true,
		MaxConnLifetime: 300 * time.Second,
	}
}
