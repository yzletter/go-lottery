package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yzletter/go-lottery/repository"
)

func Pay(ctx *gin.Context) {
	// 获取用户 id 和 gid
	uid, err := strconv.Atoi(ctx.PostForm("uid"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	gid, err := strconv.Atoi(ctx.PostForm("gid"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// 查找临时订单
	tempOrderGid := repository.GetTempOrder(uid)
	if tempOrderGid != gid {
		ctx.String(http.StatusForbidden, "您没有抢到该商品，或支付时限已过")
		return
	}

	// 支付成功，落库正式订单
	oid := repository.CreateOrder(uid, gid)
	if oid > 0 {
		repository.DeleteTempOrder(uid) // 删除临时订单

		slog.Info("用户支付成功", "uid", uid, "gid", gid)

		ctx.String(http.StatusOK, "支付成功")
		return
	} else {
		ctx.String(http.StatusInternalServerError, "系统错误，请稍后重试")
		return
	}
}

func GiveUp(ctx *gin.Context) {
	// 获取用户 id 和 gid
	uid, err := strconv.Atoi(ctx.PostForm("uid"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	gid, err := strconv.Atoi(ctx.PostForm("gid"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	repository.DeleteTempOrder(uid)

	repository.IncreaseCacheGift(gid)

	slog.Info("用户主动放弃支付", "uid", uid, "gid", gid)
}
