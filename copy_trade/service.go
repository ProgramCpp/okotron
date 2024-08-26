package copy_trade

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/pkg/errors"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/swap"
)

type Request struct {
	FromChain   string
	ToChain     string
	FromToken   string
	ToToken     string
	FromAmount  string
	FromAddress string
}

func ProcessOrder(r Request) {
	go func() {
		tradesKey := fmt.Sprintf(db.COPY_ORDER_KEY, r.FromAddress)
		users, err := db.RedisClient().LRange(context.Background(), tradesKey, 0, -1).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Printf("error fetching copy trades from redis. " + err.Error())
			return
		} else if errors.Is(err, redis.Nil) {
			// if no orders for this address,skip
			return
		}

		for _, u := range users {
			chatId, err := strconv.Atoi(u)
			if err != nil {
				log.Printf("error parsing chat id for copy order. " + err.Error())
				return
			}

			swap.SwapTokens(int64(chatId), swap.SwapRequest{
				FromToken:   r.FromToken,
				FromNetwork: r.FromChain,
				ToToken:     r.ToToken,
				ToNetwork:   r.ToChain,
				Quantity:    r.FromAmount,
			})
		}
	}()
}
