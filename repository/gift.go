package repository

import (
	"context"
	"errors"
	"fmt"
	infraMySQL "github.com/yzletter/go-lottery/infra/mysql"
	infraRedis "github.com/yzletter/go-lottery/infra/redis"
	"github.com/yzletter/go-lottery/model"
	"log/slog"
	"strconv"
)

const (
	InventoryPrefix = "gift:count:" // 方便遍历
)

// CreateCacheInventory 把 mysql 初始库存导进 redis
func CreateCacheInventory() {
	gifts := GetAllGifts()

	fmt.Println("CreateCacheInventory")

	for _, gift := range gifts {
		if gift.Count <= 0 {
			// 数据有问题
			slog.Error("Invalid Count", "name", gift.Name)
			continue
		}

		err := infraRedis.RedisClient.Set(context.Background(), InventoryPrefix+strconv.Itoa(gift.ID), gift.Count, 0).Err() // 永不过期
		if err != nil {
			slog.Error("Set Failed", "error", err)
		}
	}
}

// GetCacheInventory 从 Redis 中获得所有商品当前库存量
func GetCacheInventory() []*model.Gift {
	keys, err := infraRedis.RedisClient.Keys(context.Background(), InventoryPrefix+"*").Result()
	if err != nil {
		slog.Error("Get Failed", "error", err)
		return nil
	}

	gifts := make([]*model.Gift, 0, len(keys))
	for _, key := range keys {
		val, err := infraRedis.RedisClient.Get(context.Background(), key).Int()
		if err != nil {
			slog.Error("Get Failed", "error", err)
			continue
		}

		id, err := strconv.Atoi(key[len(InventoryPrefix):])
		if err != nil {
			continue
		}
		gifts = append(gifts, &model.Gift{
			ID:    id,
			Count: val,
		})
	}

	return gifts
}

func GetCacheGift(id int) int {
	key := InventoryPrefix + strconv.Itoa(id)
	count, err := infraRedis.RedisClient.Get(context.Background(), key).Int()
	if err != nil {
		slog.Error("Get Failed", "error", err)
		return -1
	}
	return count
}

// 库存 -1
func ReduceCacheGift(id int) error {
	key := InventoryPrefix + strconv.Itoa(id)
	count, err := infraRedis.RedisClient.Decr(context.Background(), key).Result()
	if err != nil {
		slog.Error("Get Failed", "error", err)
		return err
	} else if count < 0 {
		slog.Error("没有库存了, 仍在减库存")
		return errors.New("没有库存了, 仍在减库存")
	}
	return nil
}

// 库存 +1
func IncreaseCacheGift(id int) error {
	key := InventoryPrefix + strconv.Itoa(id)
	_, err := infraRedis.RedisClient.Incr(context.Background(), key).Result()
	if err != nil {
		slog.Error("Get Failed", "error", err)
		return err
	}
	return nil
}

// GetAllGifts 获取所有奖品
func GetAllGifts() []*model.Gift {
	var gifts []*model.Gift
	err := infraMySQL.GromDB.Model(&model.Gift{}).Select("*").Find(&gifts).Error
	if err != nil {
		slog.Error("GetAllGifts Failed", "error", err)
	}
	return gifts
}

// GetGift 根据 ID 获取奖品信息
func GetGift(id int) *model.Gift {
	var gift *model.Gift
	err := infraMySQL.GromDB.Model(&model.Gift{}).Where("id = ?", id).Find(&gift).Error
	if err != nil {
		slog.Error("GetGift Failed", "error", err)
		return nil
	}
	return gift
}
