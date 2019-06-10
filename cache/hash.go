package cache

import "github.com/garyburd/redigo/redis"

// Redis 哈希(Hash) 命令

func HashHset(key, field, value string) error {
	conn := RedisConn.Self.Get()
	if _, err := conn.Do("HSET", key, field, value); err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func HashHget(key, field string) (string, error) {
	// 不存在的key 会返回 error nil
	conn := RedisConn.Self.Get()
	value, err := redis.String(conn.Do("HGET", key, field))
	if err != nil {
		return value, err
	}
	defer conn.Close()
	return value, nil
}

func HashHexists(key, field string) (bool, error) {
	// 不存在的key 会返回 空字符串
	conn := RedisConn.Self.Get()
	exist, err := redis.Bool(conn.Do("HEXISTS", key, field))
	if err != nil {
		return exist, err
	}
	defer conn.Close()
	return exist, nil
}
