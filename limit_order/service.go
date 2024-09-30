package limit_order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	cmc "github.com/programcpp/okotron/coin_market_cap"
	"github.com/programcpp/okotron/copy_trade"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
	"github.com/programcpp/okotron/swap"
	"github.com/redis/go-redis/v9"
)

type LimitOrderRequest struct {
	ChatID int64 `json:"chat_id" redis:"/limit-order/chat-id"`
	// valid values are "buy" and "sell"
	BuyOrSell   string `json:"buy_or_sell" redis:"/limit-order/buy-or-sell"`
	FromToken   string `json:"from_token" redis:"/limit-order/from-token"`
	FromNetwork string `json:"from_network" redis:"/limit-order/from-network"`
	ToToken     string `json:"to_token" redis:"/limit-order/to-token"`
	ToNetwork   string `json:"to_network" redis:"/limit-order/to-network"`
	Quantity    string `json:"quantity" redis:"/limit-order/quantity"`
	Price       string `json:"price" redis:"/limit-order/price"`
}

func (r LimitOrderRequest) ToJson() (string, error) {
	buf := bytes.Buffer{}
	e := json.NewEncoder(&buf).Encode(r)
	return buf.String(), e
}

func (r *LimitOrderRequest) FromJson(v string) error {
	return json.NewDecoder(strings.NewReader(v)).Decode(&r)
}

// TODO. implement error channels for async process
// when the process stops silently, limit orders will no more be processed
func ProcessOrders() {
	go func() {
		for {
			time.Sleep(30 * 60 * time.Second) // for the free plan, maximum 10K calls per month. poll every 15 minutes. 4 calls per cycle

			pricesInTokens, err := cmc.PricesInTokens()
			if err != nil {
				log.Println("error fetching cmc prices in tokens")
				continue
			}

			pricesInCurrency, err := cmc.PricesInCurrency()
			if err != nil {
				log.Println("error fetching cmc prices in currency")
				continue
			}

			for token, price := range pricesInCurrency.Tokens {
				priceKey := fmt.Sprintf(db.LIMIT_ORDER_KEY, strconv.Itoa(int(price)))
				// 0: first element. -1: last element
				ordersStr, err := db.RedisClient().LRange(context.Background(), priceKey, 0, -1).Result()
				if err != nil && !errors.Is(err, redis.Nil) {
					log.Println("error fetching limit orders from redis")
					continue
				} else if errors.Is(err, redis.Nil) {
					// if no orders at this price, move to the next token price
					continue
				}

				var orders []LimitOrderRequest

				for _, os := range ordersStr {
					o := LimitOrderRequest{}
					o.FromJson(os)
					orders = append(orders, o)
				}

				// no orders at this price
				if len(orders) == 0 {
					continue
				}

				// TODO: do not check for a specific match of price, pick order within the slippage price range
				for _, o := range orders {
					if o.BuyOrSell == "buy" && o.ToToken == token || (o.BuyOrSell == "sell" && o.FromToken == token) {
						err = processOrder(o, pricesInTokens)
						if err != nil {
							log.Printf("error processing limit order. %s", err.Error())
							// do not return. the order is still in db. process next order.
							// TODO: monitor failures
						}
					}
				}
			}
		}
	}()
}

func processOrder(o LimitOrderRequest, prices cmc.PricesDataInTokens) error {
	quantity := o.Quantity
	qtyFloat, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing token quantity")
	}

	// user has entered the quantity of tokens to buy. the swap payload accepts quantity in terms of source token units
	if o.BuyOrSell == "buy" {
		quantity = fmt.Sprintf("%f", prices.Tokens[o.ToToken][o.FromToken]*qtyFloat)
	}

	err = swap.SwapTokens(o.ChatID, swap.SwapRequest{
		FromToken:   o.FromToken,
		FromNetwork: o.FromNetwork,
		ToToken:     o.ToToken,
		ToNetwork:   o.ToNetwork,
		Quantity:    quantity,
	})

	if err != nil {
		return errors.Wrap(err, "error swapping tokens in limit order")
	}

	addr, err := okto.GetAddress(o.ChatID, o.FromNetwork)
	if err != nil {
		log.Printf("error fethching addresses. %s", err.Error())
		// not affecting this transaction for copy trade failues. continue and monitor
	}

	copy_trade.ProcessOrder(copy_trade.Request{
		FromChain:   o.FromNetwork,
		FromToken:   o.FromToken,
		ToChain:     o.ToNetwork,
		ToToken:     o.ToToken,
		FromAmount:  o.Quantity,
		FromAddress: addr,
	})

	return nil
}
