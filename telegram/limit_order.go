package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/redis/go-redis/v9"
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

func LimitOrderInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// first time menu or navigated from sub menu
	// show keyboard for the next command
	var msg tgbotapi.Chattable
	if update.CallbackQuery == nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: update.FromChat().ID,
				// TODO: consolidate all telegram send messages
				ReplyMarkup: BuyOrSellKeyboard(),
			},
			Text: "buy or sell order?",
		}
	} else {
		msg = tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			"buy or sell order?", BuyOrSellKeyboard())
	}

	resp, _ := bot.Send(msg)
	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_BUY_OR_SELL,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving sub command key from-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func LimitOrderBuyOrSellInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		LimitOrderInput(bot, update)
		return
	}
	id := update.CallbackQuery.Message.MessageID
	buyOrSell := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_BUY_OR_SELL, buyOrSell).Err()
	if err != nil {
		log.Printf("error encountered when saving limit order request payload while selecting buy-or-sell. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	// TODO: move request key expiry to the primary command. its unintuitive to handle it one of the sub-commands. atleast its in the first sub-command!
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving limit order request payload while selecting buy-or-sell. %s", err.Error())
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "how would you like to pay?", tokenKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving sub command key buy-or-sell. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}


func LimitOrderFromTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		LimitOrderBuyOrSellInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_FROM_TOKEN, fromToken).Err()
	if err != nil {
		log.Printf("error encountered when saving limit order request payload while selecting from-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the source network", networkKeyboard(fromToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_FROM_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving sub command key from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func LimitOrderFromNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		LimitOrderFromTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromNetwork := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_FROM_NETWORK, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving limit order request payload while selecting from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "Which token would you like to buy?", tokenKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_TO_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func LimitOrderToTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		LimitOrderFromNetworkInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_TO_TOKEN, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving limit order request payload while selecting to-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target network", networkKeyboard(toToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_TO_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func LimitOrderToNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		LimitOrderToTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_TO_NETWORK, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving limit order request payload while selecting to-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	// TODO: the keyboard is associated with next sub-command. generalize it for all sub commands
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:", numericKeyboard())
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_QUANTITY,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key tokens command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func LimitOrderQuantiyInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		LimitOrderToNetworkInput(bot, update, false)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "enter") {
		LimitOrderCallback(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_QUANTITY)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching request payload while setting quantity. %s", res.Err())
		Send(bot, update, "something went wrong. try again.")
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_QUANTITY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while setting quantity. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	// TODO: handle error
	bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.FromChat().ID, id, "enter token quantity:"+quantity,
		numericKeyboard()))

	// the next sub command is still quantity.
	// user completes the command with this subcommand, after pressing "enter"
	// requestKey is no more accessed. will be expired by redis
}

func LimitOrder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	LimitOrderInput(bot, update)
}

func LimitOrderCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update){

}