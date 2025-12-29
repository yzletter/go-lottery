package repository

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/yzletter/go-lottery/model"
	"gorm.io/gorm"
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
func (repo *OrderRepository) CreateTempOrder(userID, giftID int) error {
	key := tempOrderPrefix + strconv.Itoa(userID)
	if err := repo.cache.Set(context.Background(), key, giftID, 0).Err(); err != nil {
		slog.Error("CreateTempOrder Failed", "error", err)
		return err
	}

	return nil
}

// GetTempOrder 查询临时订单
func (repo *OrderRepository) GetTempOrder(userID int) int {
	key := tempOrderPrefix + strconv.Itoa(userID)
	id, err := repo.cache.Get(context.Background(), key).Int()
	if err != nil {
		slog.Error("GetTempOrder Failed", "error", err)
		return 0
	}
	return id
}

// DeleteTempOrder 删除临时订单
func (repo *OrderRepository) DeleteTempOrder(userID int) int {
	key := tempOrderPrefix + strconv.Itoa(userID)

	// 返回删除订单个数
	count, err := repo.cache.Del(context.Background(), key).Result()
	if err != nil {
		slog.Error("GetTempOrder Failed", "error", err)
		return 0
	}
	return int(count)
}

func (repo *OrderRepository) CreateOrder(userID, giftID int) int {
	order := &model.Order{
		UserID: userID,
		GiftID: giftID,
	}

	err := repo.db.Model(&model.Order{}).Create(order).Error
	if err != nil {
		slog.Error("CreateOrder Failed", "error", err)
		return 0
	}

	return order.ID
}
