package repository

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/yzletter/go-lottery/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db     *gorm.DB
	client redis.UniversalClient
}

func NewOrderRepository(db *gorm.DB, client redis.UniversalClient) *OrderRepository {
	return &OrderRepository{db: db, client: client}
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
