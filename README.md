# go-whosonfirst-pool

A Redis-backed go-whosonfirst-pool LIFO pool. 

## Install

You will need to have both `Go` and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This isn't really ready to use yet.

### Simple

```
package main

import (
       "fmt"
       "github.com/whosonfirst/go-whosonfirst-pool"
       "github.com/whosonfirst/go-whosonfirst-pool-redis"
)

func main() {

     p, _ := redis.NewRedisLIFOPool("redis://localhost:6379")

     f := pool.NewIntItem(int64(123))

     p.Push(f)
     v, _ := p.Pop()

     fmt.Printf("%d", v.Int())
}
```
 
_Error handling removed for the sake of brevity._

## See also

* https://github.com/whosonfirst/go-whosonfirst-pool
* https://github.com/gomodule/redigo
* https://www.iana.org/assignments/uri-schemes/prov/redis
