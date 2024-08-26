package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/limit_order"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// wow! limit order also supports swaps together, and across networks!🔥
// limit order is essentially swap at a certain target price
type LimitOrderRequestInput struct {
	// valid values are "buy" and "sell"
	BuyOrSell   string `json:"buy_or_sell" redis:"/limit-order/buy-or-sell"`
	FromToken   string `json:"from_token" redis:"/limit-order/from-token"`
	FromNetwork string `json:"from_network" redis:"/limit-order/from-network"`
	ToToken     string `json:"to_token" redis:"/limit-order/to-token"`
	ToNetwork   string `json:"to_network" redis:"/limit-order/to-network"`
	Quantity    string `json:"quantity" redis:"/limit-order/quantity"`
	Price       string `json:"price" redis:"/limit-order/price"`
}

func (l LimitOrderRequestInput) MarshalBinary() (data []byte, err error) {
	buf := bytes.Buffer{}
	e := json.NewEncoder(&buf).Encode(l)
	return buf.Bytes(), e
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
		// TODO: edit the message to clear previous keyboards
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, resp.MessageID, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_TO_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-token command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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
		log.Printf("error encountered when saving sub command key quantity command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func LimitOrderQuantityInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data
	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	if isBack {
		LimitOrderToNetworkInput(bot, update, false)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "enter") {
		resp, _ := bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
			update.FromChat().ID, id, "enter the token limit price:",
			numericKeyboard()))

		subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
		err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_LIMIT_ORDER_CMD_PRICE,
			time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
		if err != nil {
			log.Printf("error encountered when saving sub command key price command. %s", err.Error())
			bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		}

		return
	}
	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_QUANTITY)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching request payload while setting quantity. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_QUANTITY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while setting quantity. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	msg := "enter token quantity:"
	if !strings.Contains(quantity, "back") && !strings.Contains(quantity, "enter") {
		msg += quantity
	}

	// TODO: handle error
	bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.FromChat().ID, id, msg,
		numericKeyboard()))

	// the next sub command is still quantity.
	// user completes this subcommand by pressing "enter"
}

func LimitOrderPriceInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	id := update.CallbackQuery.Message.MessageID
	price := update.CallbackQuery.Data
	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	if isBack {
		db.RedisClient().HDel(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_QUANTITY, CMD_LIMIT_ORDER_CMD_PRICE)
		LimitOrderQuantityInput(bot, update, false)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "enter") {
		LimitOrderCallback(bot, update)
		return
	}
	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_PRICE)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching request payload while setting price. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	} else if res.Err() != redis.Nil {
		price = res.Val() + price
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_LIMIT_ORDER_CMD_PRICE, price).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while setting price. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}
	// TODO: handle error
	bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.FromChat().ID, id, "enter the token limit price:"+price,
		numericKeyboard()))

	// the next sub command is still price.
	// user completes this subcommand by pressing "enter"
}

func LimitOrder(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	LimitOrderInput(bot, update)
}

func LimitOrderCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.CallbackQuery.Message.MessageID
	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	res := db.RedisClient().HGetAll(context.Background(), requestKey)
	if res.Err() != nil {
		log.Printf("error encountered when fetching limit order request payload from redis. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	var r LimitOrderRequestInput
	if err := res.Scan(&r); err != nil {
		log.Printf("error scanning limit order request payload from redis. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	limitOrderKey := fmt.Sprintf(db.LIMIT_ORDER_KEY, r.Price)
	// todo: handle error
	loReq, _ := limit_order.LimitOrderRequest{
		ChatID:      update.FromChat().ID,
		BuyOrSell:   r.BuyOrSell,
		FromToken:   r.FromToken,
		FromNetwork: r.FromNetwork,
		ToToken:     r.ToToken,
		ToNetwork:   r.ToNetwork,
		Quantity:    r.Quantity,
		Price:       r.Price,
	}.ToJson()
	// todo: handle error
	db.RedisClient().RPush(context.Background(), limitOrderKey, loReq)

	// TODO: handle error
	bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "limit order success"))
}
