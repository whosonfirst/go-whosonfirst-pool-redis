package redis

import (
	redigo "github.com/gomodule/redigo/redis"
	"github.com/whosonfirst/go-whosonfirst-pool"
)

type DeflateFunc func(pool.Item) (interface{}, error)

type InflateFunc func(interface{}, error) (pool.Item, error)

type RedisLIFOPool struct {
	pool.LIFOPool
	redis_pool *redigo.Pool
	key        string
	inflate    InflateFunc
	deflate    DeflateFunc
}

func NewRedisLIFOIntPool(dsn string, key string) (pool.LIFOPool, error) {

	deflate := func(i pool.Item) (interface{}, error) {

		return i.Int(), nil
	}

	inflate := func(rsp interface{}, err error) (pool.Item, error) {

		i, err := redigo.Int64(rsp, err)

		if err != nil {
			return nil, err
		}

		pi := pool.NewIntItem(i)
		return pi, nil
	}

	return NewRedisLIFOPool(dsn, key, deflate, inflate)
}

func NewRedisLIFOPool(dsn string, key string, deflate DeflateFunc, inflate InflateFunc) (pool.LIFOPool, error) {

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
		key:        key,
		inflate:    inflate,
		deflate:    deflate,
	}

	return &pl, nil
}

// basically the interface for pool.LIFOPool should be changed
// to expect errors all over the place but today that is not
// the case... (20181222/thisisaaronland)

// https://redis.io/commands/llen

func (pl *RedisLIFOPool) Length() int64 {

	rsp, err := pl.do("LLEN", pl.key)

	if err != nil {
		return -1
	}

	return rsp.(int64)
}

// https://redis.io/commands/rpush

func (pl *RedisLIFOPool) Push(pi pool.Item) {

	i, err := pl.deflate(pi)

	if err != nil {
		return
	}

	pl.do("LPUSH", pl.key, i)

	// error-checking?
}

// https://redis.io/commands/lpop

func (pl *RedisLIFOPool) Pop() (pool.Item, bool) {

	rsp, err := pl.do("LPOP", pl.key)

	pi, err := pl.inflate(rsp, err)

	if err != nil {
		return nil, false
	}

	return pi, true
}

func (pl *RedisLIFOPool) do(method string, args ...interface{}) (interface{}, error) {

	conn := pl.redis_pool.Get()
	defer conn.Close()

	return conn.Do(method, args...)
}
