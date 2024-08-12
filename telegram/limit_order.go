package telegram

import (

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// wow! limit order also supports swaps together, and across networks!🔥
// limit order is essentially swap at a certain target price
type LimitOrderRequestInput struct {
	BuyOrSell   bool   `json:"buy_or_sell" redis:"limit-order/buy-or-sell"`
	FromToken   string `json:"from_token" redis:"limit-order/from-token"`
	FromNetwork string `json:"from_network" redis:"limit-order/from-network"`
	ToToken     string `json:"to_token" redis:"limit-order/to-token"`
	ToNetwork   string `json:"to_network" redis:"limit-order/to-network"`
	Quantity    string `json:"quantity" redis:"limit-order/quantity"`
	Price       string `json:"price" redis:"limit-order/quantity"`
}

func LimitOrder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
}
