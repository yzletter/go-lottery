package repository

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/yzletter/go-lottery/model"
	"gorm.io/gorm"
)

type GiftRepository struct {
	db     *gorm.DB
	client redis.UniversalClient
}

func NewGiftRepository(db *gorm.DB, client redis.UniversalClient) *GiftRepository {
	return &GiftRepository{db: db, client: client}
}

// GetAllGifts 获取所有奖品
func (repo *GiftRepository) GetAllGifts() []*model.Gift {
	var gifts []*model.Gift
	err := repo.db.Model(&model.Gift{}).Select("*").Find(&gifts).Error
	if err != nil {
		slog.Error("GetAllGifts Failed", "error", err)
	}
	return gifts
}

// GetGift 根据 ID 获取奖品信息
func (repo *GiftRepository) GetGift(id int) *model.Gift {
	var gift *model.Gift
	err := repo.db.Model(&model.Gift{}).Select("id = ?", id).Find(&gift).Error
	if err != nil {
		slog.Error("GetGift Failed", "error", err)
		return nil
	}
	return gift
}
