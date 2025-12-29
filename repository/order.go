package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	infraMySQL "github.com/yzletter/go-lottery/infra/mysql"
	infraRedis "github.com/yzletter/go-lottery/infra/redis"
	"github.com/yzletter/go-lottery/model"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type OrderRepository struct {
	db    *gorm.DB
	cache redis.UniversalClient
}

const (
	tempOrderPrefix = "order:"
)

func NewOrderRepository(db *gorm.DB, client redis.UniversalClient) *OrderRepository {
	return &OrderRepository{db: db, cache: client}
}

// CreateTempOrder 创建临时订单
func CreateTempOrder(userID, giftID int) error {
	key := tempOrderPrefix + strconv.Itoa(userID)
	if err := infraRedis.RedisClient.Set(context.Background(), key, giftID, 0).Err(); err != nil {
		slog.Error("CreateTempOrder Failed", "error", err)
		return err
	}

	return nil
}

// GetTempOrder 查询临时订单
func GetTempOrder(userID int) int {
	key := tempOrderPrefix + strconv.Itoa(userID)
	id, err := infraRedis.RedisClient.Get(context.Background(), key).Int()
	if err != nil {
		slog.Error("GetTempOrder Failed", "error", err)
		return 0
	}
	return id
}

// DeleteTempOrder 删除临时订单
func DeleteTempOrder(userID int) int {
	key := tempOrderPrefix + strconv.Itoa(userID)

	// 返回删除订单个数
	count, err := infraRedis.RedisClient.Del(context.Background(), key).Result()
	if err != nil {
		slog.Error("GetTempOrder Failed", "error", err)
		return 0
	}
	return int(count)
}

func CreateOrder(userID, giftID int) int {
	order := &model.Order{
		UserID: userID,
		GiftID: giftID,
	}

	err := infraMySQL.GromDB.Model(&model.Order{}).Create(order).Error
	if err != nil {
		slog.Error("CreateOrder Failed", "error", err)
		return 0
	}

	return order.ID
}
