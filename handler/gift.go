package handler

import (
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yzletter/go-lottery/model"
	"github.com/yzletter/go-lottery/mq"
	"github.com/yzletter/go-lottery/repository"
)

const (
	PayDelay = 600
)

// 获取所有奖品信息
func GetAllGifts(ctx *gin.Context) {
	gifts := repository.GetAllGifts()
	for _, gift := range gifts {
		// 去掉敏感信息
		gift.Count = 1
	}

	ctx.JSON(200, gifts)
}

func Lottery(ctx *gin.Context) {
	// 尝试十次
	for try := 1; try <= 10; try++ {
		// 获取所有库存
		gifts := repository.GetCacheInventory()

		ids := make([]int, 0, len(gifts))
		prob := make([]float64, 0, len(gifts))

		for _, gift := range gifts {
			if gift.Count > 0 {
				ids = append(ids, gift.ID)
				prob = append(prob, float64(gift.Count))
			}
		}

		// 没获取到库存
		if len(gifts) == 0 {
			ctx.String(200, strconv.Itoa(0)) // 0 表示抽完了
			return
		}

		index := lottery(prob)
		gid := ids[index]

		err := repository.ReduceCacheGift(gid)
		if err != nil {
			slog.Error("减库存失败", "error", err)
			continue
		}

		uid := 1                        // 没做登录系统，把用户id写死为1
		gift := repository.GetGift(gid) // 获取物品详情
		if gift == nil {
			slog.Error("找不到奖品", "gid", gid)
			continue
		}

		// 创建临时订单
		repository.CreateTempOrder(uid, gid)

		// 发送延迟消息
		mq.Send(&model.Order{UserID: uid, GiftID: gid}, PayDelay)

		slog.Info("抽中奖品", "用户", uid, "奖品", gid)

		// 先设置Cookie
		ctx.SetCookie("name", gift.Name, PayDelay, "/", "localhost", false, false)                 // 抢中的商品名称
		ctx.SetCookie("price", strconv.Itoa(gift.Price), PayDelay, "/", "localhost", false, false) // 商品价格
		ctx.SetCookie("uid", strconv.Itoa(uid), PayDelay, "/", "localhost", false, false)          // 用户id
		ctx.SetCookie("gid", strconv.Itoa(gid), PayDelay, "/", "localhost", false, false)          // 商品id

		// 再设置body
		ctx.String(http.StatusOK, strconv.Itoa(gid)) // 减库存成功后才给前端返回奖品ID

		return
	}
}

func lottery(probs []float64) int {
	if len(probs) == 0 {
		return -1
	}

	sum := 0.0
	acc := make([]float64, 0, len(probs))
	for _, prob := range probs {
		sum += prob
		acc = append(acc, sum)
	}

	// 获取(0,sum] 随机数
	x := rand.Float64() * sum
	l, r := 0, len(probs)-1
	for l < r {
		mid := (l + r) / 2
		if acc[mid] < x {
			l = mid + 1
		} else {
			r = mid
		}
	}

	return l
}
