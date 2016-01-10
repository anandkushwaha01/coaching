package common

// QZ: based on golang-relyq code

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

type RedisStorage struct {
	pool   *redis.Pool
	prefix string
}

func InitRedis(address string) (*redis.Pool, error) {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	if pool == nil {
		return nil, errors.New("pool not found")
	}
	return pool, nil
}
func GetRedisPool() *redis.Pool {
	return pool
}
func NewRedisStorage(pool *redis.Pool, prefix string, delim string) *RedisStorage {
	return &RedisStorage{
		pool:   pool,
		prefix: prefix + delim + "oa" + delim,
	}
}

func (rs *RedisStorage) Get(id []byte, obj interface{}) error {
	val, err := redis.Bytes(rs.do("GET", rs.prefixed(id)))

	if err != nil {
		return err
	}

	return json.Unmarshal(val, obj)
}

func (rs *RedisStorage) Set(obj interface{}, id []byte) error {
	val, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = rs.do("SET", rs.prefixed(id), val)
	return err
}

func (rs *RedisStorage) Del(id []byte) error {
	_, err := rs.do("DEL", rs.prefixed(id))
	return err
}

func (rs *RedisStorage) prefixed(id []byte) []byte {
	return append([]byte(rs.prefix), id...)
}

func (rs *RedisStorage) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := rs.pool.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

func (rs *RedisStorage) SetEx(obj interface{}, id []byte, expireSeconds int) error {
	val, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = rs.do("SETEX", rs.prefixed(id), expireSeconds, val)
	//_, err = rs.do("EXPIRE", rs.prefixed(id), 3600)
	return err
}
