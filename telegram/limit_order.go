package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/spf13/viper"
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
	// first time menu or navigated from sub menu
	// show keyboard for the next command
	var msg tgbotapi.Chattable
	if update.CallbackQuery == nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: update.FromChat().ID,
				// TODO: consolidate all telegram send messages
				ReplyMarkup: tokenKeyboard(),
			},
			Text: "select the source token",
		}
	} else {
		msg = tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			"select the source token", tokenKeyboard())
	}

	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), messageKey, CMD_ANY_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key from-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}
