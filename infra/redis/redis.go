package infra

import (
	"context"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/yzletter/go-lottery/infra/viper"
)

var (
	RedisClient *redis.Client
	redisOnce   sync.Once
)

// Init 连接到 Redis 数据库, 生成一个 *redis.Client 赋给全局数据库变量 RedisClient
func Init(confDir, confFileName, confFileType string) redis.UniversalClient {
	// 初始化 Viper 进行配置读取
	viper := viper.InitViper(confDir, confFileName, confFileType)
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	db := viper.GetInt("redis.db")

	redisAddr := host + ":" + port // 拼接地址
	redisOption := &redis.Options{
		Addr: redisAddr,
		DB:   db,
	}

	// 连接到数据库
	redisOnce.Do(func() {
		RedisClient = redis.NewClient(redisOption)
	})

	// 尝试 ping 通
	if err := RedisClient.Ping(context.Background()).Err(); err != nil { // 须加上.Err(), 否则会报 ping 通错
		slog.Error("connect to Redis failed", "error", err)
		panic(err)
	} else {
		slog.Info("connect to Redis succeed")
	}

	return RedisClient
}

// Ping ping 一下数据库 保持连接
func Ping() {
	if RedisClient != nil {
		err := RedisClient.Ping(context.Background()).Err()
		if err != nil {
			slog.Info("ping RedisClient failed")
			return
		}
		slog.Info("ping RedisClient succeed")
		return
	}
}

func Close() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			slog.Info("close RedisClient failed")
			return
		}
		slog.Info("close RedisClient succeed")
		return
	}
}
