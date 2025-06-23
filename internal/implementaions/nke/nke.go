package nke

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/dto/message"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/services"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

type NKE struct {
	SqlStore    sqlc.IStore
	RDb         *redis.Client
	RedisClient services.IRedis
}

func (n *NKE) RunSubscribingToExpiredEvents(ctx context.Context) {
	pubsub := n.RDb.PSubscribe(ctx, "__keyevent@0__:expired")
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopped subscribing due to context cancel")
			return
		case msg := <-ch:
			key := msg.Payload
			if strings.HasPrefix(key, global.ORDER) {
				orderId := strings.TrimPrefix(key, global.ORDER)

				var orderPayloadRedisKey = fmt.Sprintf("%s%s", global.ORDER_PAYLOAD, orderId)
				var subOrder message.PayloadSubOrderData

				err := n.RedisClient.Get(orderPayloadRedisKey, &subOrder)
				if err != nil {
					if err.Error() == fmt.Sprintf("key %s does not exist", orderPayloadRedisKey) {
						log.Println(err.Error())
						continue
					}

					log.Printf("failed to retrieve payload from Redis with key %s: %v\n", orderPayloadRedisKey, err)
					continue
				}

				err = n.SqlStore.UpdateSeatStatusTran(ctx, subOrder, global.ORDER_FAILED)
				if err != nil {
					log.Printf("failed to update order status to ORDER_FAILED for orderId=%s: %v\n", orderId, err)
					continue
				}

				log.Printf("Successfully handled expired orderId=%s and updated status to ORDER_FAILED.\n", orderId)
			}
		}
	}
}

func NewNKE(sqlStore sqlc.IStore,
	redisClient services.IRedis) *NKE {
	return &NKE{
		RDb:         global.RDb,
		SqlStore:    sqlStore,
		RedisClient: redisClient,
	}
}
