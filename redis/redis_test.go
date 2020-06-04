package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestRedisConn(t *testing.T) {
	conn, err := RedisConn()
	if err != nil {
		t.Fatal("redis connect failed" + err.Error())
	}
	fmt.Println(conn)
}

func TestNewPool(t *testing.T) {
	pool := NewPool()
	defer pool.Close()
	fmt.Println(pool)

	for i := 0; i <= 3; i++ {
		go func() {
			c := pool.Get()
			defer c.Close()

			_, err := c.Do("MSET", "testname", "baidu", "testurl", "www.baidu.com")
			if err != nil {
				t.Fatal("mset value error : " + err.Error())
			}
			if mget, err := redis.Strings(c.Do("MGET", "testname", "testurl")); err == nil {
				for _, v := range mget {
					fmt.Println(v)
				}
			}
		}()
	}
	time.Sleep(2 * time.Second)
}
