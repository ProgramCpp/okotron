package limitorder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	_ "time"

	cmc "github.com/programcpp/okotron/coin_market_cap"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/telegram"
	"github.com/redis/go-redis/v9"
)

// TODO. implement error channels for async process
// when the process stops silently, limit orders will no more be processed
func ProcessOrders() {
	go func() {
		for {
			// time.Sleep(5 * 60 * time.Second) // for the free plan, maximum 10K calls per month. poll every 5 minutes

			prices, err := cmc.Prices()
			if err != nil {
				log.Println("error fetching cmc prices")
				return
			}

			for token, price := range prices.Tokens {

				priceKey := fmt.Sprintf(db.LIMIT_ORDER_KEY, strconv.Itoa(int(price)))
				// 0: first element. -1: last element
				ordersResult := db.RedisClient().LRange(context.Background(), priceKey, 0, -1)
				if ordersResult.Err() != nil && !errors.Is(ordersResult.Err(), redis.Nil) {
					log.Println("error fetching limit orders from redis")
					return
				} else if errors.Is(ordersResult.Err(), redis.Nil) {
					// if no orders at this price, move to the next token price
					continue
				}

				orders := []telegram.LimitOrderRequestInput{}
				ordersResult.ScanSlice(orders)

				for _, o := range orders {
					if o.BuyOrSell == "buy" && o.ToToken == token {
						processOrder(o)
					} else if o.BuyOrSell == "sell" && o.FromToken == token {
						processOrder(o)
					}
				}
			}
		}
	}()
}


func processOrder(o telegram.LimitOrderRequestInput){

}