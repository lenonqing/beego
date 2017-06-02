package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	defaultMaxActiveSize = 500
	defaultMaxIdleSize   = 100
	defaultIdleTimeout   = 600 * time.Second
)

type redisPool interface {
	ActiveCount() int
	Close() error
	Get() redis.Conn
}

// RedisDB 数据库实例结构.
type RedisDB struct {
	pool redisPool
}

func (db *RedisDB) ping() (bool, error) {
	conn := db.pool.Get()
	defer conn.Close()
	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

func dial(address, auth string) (redis.Conn, error) {
	c, err := redis.Dial(`tcp`, address)
	if err != nil {
		return c, err
	}
	if auth != "" {
		if _, err := c.Do("AUTH", auth); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

// CreateRedisDB 创建redis数据库连接实例.
func CreateRedisDB(address, auth string) (db *RedisDB, err error) {
	pool := &redis.Pool{
		MaxActive:   defaultMaxActiveSize,
		MaxIdle:     defaultMaxIdleSize,
		IdleTimeout: defaultIdleTimeout,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial(address, auth)
		},
	}

	db = &RedisDB{
		pool: pool,
	}
	_, err = db.ping()
	return
}

///////////////////接口///////////////////

// Get 获取key对应的value
func (db *RedisDB) Get(key string) (content string, err error) {
	conn := db.pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	content, err = redis.String(conn.Do("GET", key))
	return
}

// Set 设置key对应的值
func (db *RedisDB) Set(key string, content string, age uint) (err error) {
	conn := db.pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	if age > 0 {
		_, err = conn.Do("SET", key, content, "EX", age)
	} else {
		_, err = conn.Do("SET", key, content)
	}
	return
}

// HMSet 设置key对应的值
func (db *RedisDB) HMSet(key string, json map[string]interface{}) (err error) {
	conn := db.pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("HMSet", redis.Args{}.Add(key).AddFlat(json)...)
	return
}

// HSet 设置key对应的值
func (db *RedisDB) HSet(key string, field string, value string) (err error) {
	conn := db.pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("HSet", key, field, value)
	return
}

// HGetAll 获取对应的值
func (db *RedisDB) HGetAll(key string) (json map[string]string, err error) {
	conn := db.pool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	json, err = redis.StringMap(conn.Do("HGETALL", key))
	return
}

// SAdd 添加到set中
func (db *RedisDB) SAdd(key string, members []string) (err error) {
	conn := db.pool.Get()
	defer conn.Close()
	args := redis.Args{}
	if _, err = conn.Do("SADD", args.Add(key).AddFlat(members)...); err != nil {
		return err
	}
	return nil
}

// SIsMember member是否在set中
func (db *RedisDB) SIsMember(key string, member string) (isMember bool, err error) {
	conn := db.pool.Get()
	defer conn.Close()
	isMember, err = redis.Bool(conn.Do("SISMEMBER", key, member))
	return
}

// Expire 设置过期时间
func (db *RedisDB) Expire(key string, seconds int64) (err error) {
	conn := db.pool.Get()
	defer conn.Close()
	if _, err = conn.Do("EXPIRE", key, seconds); err != nil {
		return err
	}
	return nil
}
