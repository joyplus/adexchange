package lib

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     500,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if len(password) > 0 {
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

var PoolQueue *redis.Pool
var PoolCache *redis.Pool

func SetQueuePool(pool *redis.Pool) {
	PoolQueue = pool
}

func SetCachePool(pool *redis.Pool) {
	PoolCache = pool
}

func GetQueuePool() *redis.Pool {
	return PoolQueue
}

func GetCachePool() *redis.Pool {
	return PoolCache
}

//var (
//	redisServer   = flag.String("redisServer", ":6379", "")
//	redisPassword = flag.String("redisPassword", "", "")
//)
