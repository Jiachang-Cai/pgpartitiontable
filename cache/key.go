package cache

// Redis 键(key) 命令

func KeyDel(key string) error {
	conn := RedisConn.Self.Get()
	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func KeyExpireAt(key string, timestamp int64) error {
	conn := RedisConn.Self.Get()
	defer conn.Close()
	if _, err := conn.Do("EXPIREAT", key, timestamp); err != nil {
		return err
	}
	return nil
}

func KeyExpire(key string, seconds int64) error {
	conn := RedisConn.Self.Get()
	defer conn.Close()
	if _, err := conn.Do("EXPIRE", key, seconds); err != nil {
		return err
	}
	return nil
}
