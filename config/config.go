package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var QConfig *Config

type Config struct {
	Redis  RedisConfig
	Log    LogConfig
	DelayQ DelayQConfig
	App    ApplicationConfig
}

type RedisConfig struct {
	IP              string
	Port            int
	Db              int
	Password        string
	MaxIdle         int
	MaxActive       int
	IdleTimeout     int
	MaxConnLifetime int
	ReadTimeout     int
	WriteTimeout    int
	KeepAlive       int
	ConnTimeout     int
}

type LogConfig struct {
	ErrorLog  string
	AccessLog string
	LogLevel  string
	LogEncode string
}

// Delay Queue Config
type DelayQConfig struct {
	ReadyQueue string
	JobPoll    string
	DelayQueue string
}

type ApplicationConfig struct {
	Port int
	IP   string
}

func InitConfig() {
	QConfig, _ = LoadConfig()
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(getConfigDir())
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("viper read config failed")
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("unmarshal config failed")
	}
	return &cfg, nil
}

func GetRedisConfig() *RedisConfig {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println("read config failed")
		return nil
	}
	return &config.Redis
}

func getConfigDir() string {
	dir, err := filepath.Abs(filepath.Dir("./"))
	if err != nil {
		panic("read config dir error: " + err.Error())
	}
	dir = strings.Replace(dir, "\\", "/", -1)
	return strings.Join([]string{dir, "config"}, "/")
}

func NewRedisConfig(config *Config) *RedisConfig {
	getredis := config.Redis
	return &RedisConfig{
		IP:              getredis.IP,
		Port:            getredis.Port,
		Db:              getredis.Db,
		Password:        getredis.Password,
		MaxIdle:         getredis.MaxIdle,
		MaxActive:       getredis.MaxActive,
		IdleTimeout:     getredis.IdleTimeout,
		MaxConnLifetime: getredis.MaxConnLifetime,
		ReadTimeout:     getredis.ReadTimeout,
		WriteTimeout:    getredis.WriteTimeout,
		KeepAlive:       getredis.KeepAlive,
		ConnTimeout:     getredis.ConnTimeout,
	}
}
