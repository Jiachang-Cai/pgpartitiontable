package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-xweb/log"
	"github.com/spf13/viper"
)

type DataCache struct {
	Self *redis.Pool
}

var RedisConn *DataCache

// newPool New redis pool.
func newPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server, redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					log.Errorf("occur error at newPool: %v\n", err)
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

func (cache *DataCache) Init() {
	RedisConn = &DataCache{
		Self: newPool(viper.GetString("redisservers.cajiancache.addr"),
			viper.GetString("redisservers.cajiancache.password"), viper.GetInt("redisservers.cajiancache.db")),
	}
}

