package redis

import (
	redigo "github.com/gomodule/redigo/redis"
	"github.com/whosonfirst/go-whosonfirst-pool"
	_ "log"
)

type RedisLIFOPool struct {
	pool.LIFOPool
	redis_pool *redigo.Pool
	key        string
}

func NewRedisLIFOPool(dsn string) (pool.LIFOPool, error) {

	redis_pool := &redigo.Pool{
		MaxActive: 1000,
		Dial: func() (redigo.Conn, error) {

			// https://www.iana.org/assignments/uri-schemes/prov/redis

			c, err := redigo.DialURL(dsn)

			if err != nil {
				return nil, err
			}

			return c, err
		},
	}

	pl := RedisLIFOPool{
		redis_pool: redis_pool,
		key:        "debug",
	}

	return &pl, nil
}

// https://redis.io/commands/llen

func (pl *RedisLIFOPool) Length() int64 {

	rsp, err := pl.do("LLEN", pl.key)

	if err != nil {
		return -1
	}

	return rsp.(int64)
}

// https://redis.io/commands/rpush

func (pl *RedisLIFOPool) Push(i pool.Item) {

	pl.do("LPUSH", pl.key, i.Int())

	// error-checking?
}

// https://redis.io/commands/lpop

func (pl *RedisLIFOPool) Pop() (pool.Item, bool) {

	rsp, err := pl.do("LPOP", pl.key)

	if err != nil {
		return nil, false
	}

	i, err := redigo.Int64(rsp, err)

	if err != nil {
		return nil, false
	}

	pi := pool.NewIntItem(i)
	return pi, true
}

func (pl *RedisLIFOPool) do(method string, args ...interface{}) (interface{}, error) {

	conn := pl.redis_pool.Get()
	defer conn.Close()

	return conn.Do(method, args...)
}
