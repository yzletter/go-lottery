package mq

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/bytedance/sonic"
	"github.com/yzletter/go-lottery/model"
)

const (
	END_POINT = "localhost:8081"
	// ./mqadmin.cmd updateTopic -n localhost:9876 -c DefaultCluster -t CANCEL_ORDER -a +message.type=DELAY
	// ./mqadmin.cmd deleteTopic -n localhost:9876 -c DefaultCluster -t CANCEL_ORDER
	TOPIC = "CANCEL_ORDER"
	// ./mqadmin.cmd updateSubGroup -n localhost:9876 -c DefaultCluster -g lottery
	CONSUMER_GROUP = "go_lottery"
)

func init() {
	// 将 rocketmq 的日志输出到控制台
	os.Setenv(rmq_client.ENABLE_CONSOLE_APPENDER, "true")
	rmq_client.ResetLogger()
}

var (
	producer rmq_client.Producer
	ponce    sync.Once
)

func GetProducer() rmq_client.Producer {
	ponce.Do(func() {
		// 初始化过
		if producer != nil {
			return
		}

		// 未初始化过
		var err error
		producer, err = rmq_client.NewProducer(
			&rmq_client.Config{
				Endpoint:      END_POINT,
				NameSpace:     "",
				ConsumerGroup: "",
				Credentials:   &credentials.SessionCredentials{},
			},
			rmq_client.WithTopics(TOPIC),
		)
		if err != nil {
			log.Fatal(err)
		}

		err = producer.Start()
		if err != nil {
			log.Fatal(err)
		}
	})

	return producer
}

func StopProducer() {
	if producer != nil {
		producer.GracefulStop()
		slog.Info("stop producer")
	}
}

func Send(msg *model.Order, delay int) error {
	producer := GetProducer()

	ctx := context.Background()
	body, err := sonic.Marshal(msg)
	if err != nil {
		return err
	}

	message := &rmq_client.Message{
		Topic: TOPIC,
		Body:  body,
	}
	message.SetDelayTimestamp(time.Now().Add(time.Second * time.Duration(delay)))

	_, err = producer.Send(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
