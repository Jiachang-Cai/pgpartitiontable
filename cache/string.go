package cache

import "github.com/garyburd/redigo/redis"

// Redis 字符串(String) 命令

func StringSet(key string, value interface{}) error {
	conn := RedisConn.Self.Get()
	if _, err := conn.Do("SET", key, value); err != nil {
		return err
	}

	defer conn.Close()
	return nil
}

func StringGet(key string) (string, error) {
	conn := RedisConn.Self.Get()
	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return value, err
	}
	defer conn.Close()
	return value, nil
}

func StringIncr(key string) (int, error) {
	conn := RedisConn.Self.Get()
	value, err := redis.Int(conn.Do("INCR", key))
	if err != nil {
		return value, err
	}
	defer conn.Close()
	return value, nil
}

func StringSetExpire(key, value string, seconds int64) error {
	conn := RedisConn.Self.Get()
	if _, err := conn.Do("SET", key, value, "EX", seconds); err != nil {
		return err
	}
	defer conn.Close()
	return nil
}


