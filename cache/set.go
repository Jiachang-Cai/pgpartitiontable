package cache

import "github.com/garyburd/redigo/redis"

// Redis 集合(Set) 命令

// 移除集合中一个或多个成员
func SetSrem(key string, arr ...string) error {
	args := []interface{}{key}
	for _, v := range arr {
		args = append(args, v)
	}
	conn := RedisConn.Self.Get()
	if _, err := conn.Do("SREM", args...); err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// 向集合添加一个或多个成员
func SetSadd(key string, arr ...string) error {
	args := []interface{}{key}
	for _, v := range arr {
		args = append(args, v)
	}
	conn := RedisConn.Self.Get()
	if _, err := conn.Do("SADD", args...); err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// 命令判断成员元素是否是集合的成员
func SetSismember(key, value string) (bool, error) {
	conn := RedisConn.Self.Get()
	exist, err := redis.Bool(conn.Do("SISMEMBER", key, value))
	if err != nil {
		return exist, err
	}
	defer conn.Close()
	return exist, nil
}
