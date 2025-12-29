package main

import (
	infraMySQL "github.com/yzletter/go-lottery/infra/mysql"
	infraRedis "github.com/yzletter/go-lottery/infra/redis"
	"github.com/yzletter/go-lottery/infra/slog"
	"github.com/yzletter/go-lottery/infra/viper"
)

func main() {
	slog.InitSlog("./logs/go_lottery.log")                          // 初始化 slog
	GormDB := infraMySQL.Init("./conf", "db", viper.YAML, "./logs") // 注册 MySQL
	RedisClient := infraRedis.Init("./conf", "cache", viper.YAML)   // 初始化 Redis

}
