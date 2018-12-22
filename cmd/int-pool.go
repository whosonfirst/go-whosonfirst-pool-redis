package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-pool"
	"github.com/whosonfirst/go-whosonfirst-pool-redis"
	"log"
)

func main() {

	var dsn = flag.String("dsn", "redis://localhost:6379", "The data source name (dsn) for connecting to Redis.")

	flag.Parse()

	p, err := redis.NewRedisLIFOPool(*dsn)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("LEN", p.Length())

	f1 := pool.NewIntItem(int64(123))
	f2 := pool.NewIntItem(int64(456))

	p.Push(f1)
	p.Push(f2)

	v, ok := p.Pop()

	if !ok {
		log.Fatal("Did not pop")
	}

	fmt.Printf("%d", v.Int())
}