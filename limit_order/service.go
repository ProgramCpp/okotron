package limitorder

import (
	"log"
	_ "time"

	cmc "github.com/programcpp/okotron/coin_market_cap"
)

// TODO. implement error channels for async process
func ProcessOrders() {
	go func() {
		for {
			// time.Sleep(5 * 60 * time.Second) // for the free plan, maximum 10K calls per month. poll every 5 minutes

			prices, err := cmc.Prices()
			if err != nil {
				log.Println("error fetching cmc prices")
				return
			}

			_ = prices
		}
	}()
}
