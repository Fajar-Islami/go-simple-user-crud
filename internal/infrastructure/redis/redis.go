package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	redisRepo "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/redis"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
	"github.com/go-redis/redis/v8"
)

type RedisConf struct {
	Password          string `env:"redis_password"`
	Port              int    `env:"redis_port"`
	Host              string `env:"redis_host"`
	DB                int    `env:"redis_db"`
	DefaultDB         int    `env:"redis_Defaultdb"`
	RedisMinIdleConns int    `env:"redis_MinIdleConns"`
	RedisPoolSize     int    `env:"redis_PoolSize"`
	RedisPoolTimeout  int    `env:"redis_PoolTimeout"`
	RedisTTL          int    `env:"redis_ttl"`
	TLSConfig         bool   `env:"redis_TLSConfig"`
}

const currentfilepath = "internal/infrastructure/redis/redis.go"

func NewRedisClient() *redis.Client {
	var redisConfig = RedisConf{
		Password:          utils.EnvString("redis_password"),
		Port:              utils.EnvInt("redis_port"),
		Host:              utils.EnvString("redis_host"),
		DB:                utils.EnvInt("redis_db"),
		DefaultDB:         utils.EnvInt("redis_Defaultdb"),
		RedisMinIdleConns: utils.EnvInt("redis_MinIdleConns"),
		RedisPoolSize:     utils.EnvInt("redis_PoolSize"),
		RedisPoolTimeout:  utils.EnvInt("redis_PoolTimeout"),
		RedisTTL:          utils.EnvInt("redis_ttl"),
		TLSConfig:         utils.EnvBool("redis_TLSConfig"),
	}

	redisRepo.RedisTTL = time.Duration(redisConfig.RedisTTL * int(time.Second))

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		MinIdleConns: redisConfig.RedisMinIdleConns,
		PoolSize:     redisConfig.RedisPoolSize,
		PoolTimeout:  time.Duration(redisConfig.RedisPoolTimeout) * time.Second,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
	})

	if redisConfig.TLSConfig {
		client.Options().TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelPanic, "", fmt.Errorf("Cannot conenct to redis : %s", err.Error()))
		panic(err)
	}
	helper.Logger(currentfilepath, helper.LoggerLevelInfo, fmt.Sprintf("Redis ping : %s", pong), nil)

	return client
}
