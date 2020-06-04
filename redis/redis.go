package redis

import (
	"fmt"
	"time"

	"github.com/jiujuan/delayq/config"

	"github.com/gomodule/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

func InitRedis() {
	RedisPool = NewPool()
}

//redisConn 连接redis
func RedisConn() (redis.Conn, error) {
	conf := config.QConfig.Redis
	host := fmt.Sprintf("%s:%d", conf.IP, conf.Port)
	conn, err := redis.Dial(
		"tcp",
		host,
		redis.DialConnectTimeout(parseTime(conf.ConnTimeout)),
		redis.DialReadTimeout(parseTime(conf.ReadTimeout)),
		redis.DialWriteTimeout(parseTime(conf.WriteTimeout)),
		redis.DialPassword(conf.Password),
		redis.DialKeepAlive(parseTime(conf.KeepAlive)),
	)
	return conn, err
}

//NewPool 构造连接池
func NewPool() *redis.Pool {
	conf := config.QConfig.Redis
	return &redis.Pool{
		MaxIdle:         conf.MaxIdle,
		MaxActive:       conf.MaxActive,
		IdleTimeout:     parseTime(conf.IdleTimeout),
		MaxConnLifetime: parseTime(conf.MaxConnLifetime),
		Dial:            func() (redis.Conn, error) { return RedisConn() },
	}
}

type RedisVal struct {
	Key    string        `json:"key"`
	Field  string        `json:"field"`
	Value  interface{}   `json:"value"`
	Values []interface{} `json:"values"`
}

// Execute 执行redis命令
func Do(command string, args ...interface{}) (interface{}, error) {
	pool := RedisPool.Get()
	defer pool.Close()

	return pool.Do(command, args...)
}

func HSET(args ...interface{}) (string, error) {
	s, err := redis.String(Do("HSET", args...))
	return s, err
}

func HGET(args ...interface{}) ([]string, error) {
	s, err := redis.Strings(Do("HGET", args...))
	return s, err
}

func HEXISTS(args ...interface{}) (string, error) {
	s, err := redis.String(Do("HEXISTS", args))
	return s, err
}

func SET(args ...interface{}) (string, error) {
	s, err := redis.String(Do("SET", args...))
	return s, err
}

func GET(args ...interface{}) (string, error) {
	s, err := redis.String(Do("GET", args...))
	return s, err
}

func DEL(args ...interface{}) error {
	_, err := redis.String(Do("DEL", args...))
	return err
}

func ZADD(args ...interface{}) (string, error) {
	s, err := redis.String(Do("ZADD", args...))
	return s, err
}

func ZRANGEBYSCORE(args ...interface{}) ([]string, error) {
	s, err := redis.Strings(Do("ZRANGEBYSCORE", args...))
	return s, err
}

func ZREM(args ...interface{}) error {
	_, err := Do("ZREM", args...)
	return err
}

func LPUSH(args ...interface{}) error {
	_, err := Do("LPUSH", args...)
	return err
}

func BRPOP(args ...interface{}) (string, error) {
	s, err := redis.String(Do("BRPOP", args...))
	return s, err
}

func parseTime(num int) time.Duration {
	return time.Duration(num) * time.Second
}

func getRedisConfig() *config.RedisConfig {
	getconfig, err := config.LoadConfig()
	if err != nil {
		panic("load config error")
	}
	return &getconfig.Redis
}
