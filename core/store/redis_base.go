// @Title
// @Description
// @Author  Wangwengang  2023/12/23 12:21
// @Update  Wangwengang  2023/12/23 12:21
package store

import (
	"context"
	"fmt"
	"hash/crc32"
	"runtime"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog"
	"go.uber.org/zap"
)

type RedisBase struct {
	RedisCli *redis.Client
}

var redisBaseInstance *RedisBase

func setRedisBase(redisBase *RedisBase) {
	redisBaseInstance = redisBase
}

func RedisIns() *RedisBase {
	if redisBaseInstance == nil {
		slog.Ins().Errorf("redisBase is nil")
		return nil
	}
	return redisBaseInstance
}

func NewCache(config sconfig.Redis) *RedisBase {
	numCPU := runtime.NumCPU()
	redisCli := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   3,
		PoolSize:     numCPU * 512,
		MinIdleConns: numCPU * 8,
	})

	pong, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis connect ping failed, err:", zap.Error(err))
		panic(err)
	} else {
		fmt.Println("redis connect ping response:", zap.String("pong", pong))
		rb := &RedisBase{
			RedisCli: redisCli,
		}
		setRedisBase(rb)
		return rb
	}
}

func (r *RedisBase) GetHash(k int64) int64 {
	s := strconv.FormatInt(k, 10)
	v := int64(crc32.ChecksumIEEE([]byte(s)))
	if v < 0 {
		v = -v
	}

	v = v % 100
	return v
}
