package mq

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/bytedance/sonic"
	"github.com/yzletter/go-lottery/model"
	"github.com/yzletter/go-lottery/repository"
)

func init() {
	// 将 rocketmq 的日志输出到控制台
	os.Setenv(rmq_client.ENABLE_CONSOLE_APPENDER, "true")
	rmq_client.ResetLogger()
}

var (
	simpleConsumer rmq_client.SimpleConsumer
	once           sync.Once
)

func GetConsumer() rmq_client.SimpleConsumer {
	once.Do(func() {
		// 初始化过
		if producer != nil {
			return
		}

		// 未初始化过
		var err error
		simpleConsumer, err = rmq_client.NewSimpleConsumer(
			&rmq_client.Config{
				Endpoint:      END_POINT, // Proxy 地址
				Credentials:   &credentials.SessionCredentials{},
				ConsumerGroup: CONSUMER_GROUP, // 消费方需要指定组
				NameSpace:     "",
			},
			rmq_client.WithSimpleAwaitDuration(5*time.Second),
			rmq_client.WithSimpleSubscriptionExpressions(map[string]*rmq_client.FilterExpression{
				TOPIC: rmq_client.SUB_ALL, // 订阅该 Topic 下所有 Tag
			}),
		)
		if err != nil {
			log.Fatal(err)
		}

		err = simpleConsumer.Start()
		if err != nil {
			log.Fatal(err)
		}
	})

	return simpleConsumer
}

func StopConsumer() {
	if simpleConsumer != nil {
		simpleConsumer.GracefulStop()
		slog.Info("stop consumer")
	}
}

func Consume() {
	consumer := GetConsumer()
	ctx := context.Background()
	for {
		messages, err := consumer.Receive(ctx, 1, 10*time.Second) // 一批一条
		if err != nil {
			// 判断是否 broker 里暂时没有数据, 40401
			var e *rmq_client.ErrRpcStatus
			if errors.As(err, &e) && e.Code != 40401 {
				log.Printf("Receive Message Failed, Code %d, Error %s\n", e.Code, e.Message)
			}
			continue
		}

		for _, message := range messages {
			var order model.Order
			err := sonic.Unmarshal(message.GetBody(), &order)
			if err != nil {
				continue
			}
			gid := repository.GetTempOrder(order.UserID)
			if gid == order.GiftID {
				// 支付超时，删除临时订单，增加库存
				repository.DeleteTempOrder(order.UserID)
				repository.IncreaseCacheGift(gid)
			}
			consumer.Ack(ctx, message)
		}
	}
}
