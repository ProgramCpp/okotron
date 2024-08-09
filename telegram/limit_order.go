package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// wow! limit order also supports swaps, and across networks!🔥
type LimitOrderRequestInput struct {
	FromToken   string `json:"from_token" redis:"limit-order/from-token"`
	FromNetwork string `json:"from_network" redis:"limit-order/from-network"`
	ToToken     string `json:"to_token" redis:"limit-order/to-token"`
	ToNetwork   string `json:"to_network" redis:"limit-order/to-network"`
	Quantity    string `json:"quantity" redis:"limit-order/quantity"`
	Price       string `json:"price" redis:"limit-order/quantity"`
}

func LimitOrder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

}
