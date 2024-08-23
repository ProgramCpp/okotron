package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/swap"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// redis tags are defined as constants. keep in SYNC. the keys are saved with COMMAND Names and deserialized with redis tags!
type SwapRequestInput struct {
	FromToken   string `json:"from_token" redis:"/swap/from-token"`
	FromNetwork string `json:"from_network" redis:"/swap/from-network"`
	ToToken     string `json:"to_token" redis:"/swap/to-token"`
	ToNetwork   string `json:"to_network" redis:"/swap/to-network"`
	Quantity    string `json:"quantity" redis:"/swap/quantity"` // use okto terms - quantity. not lifi terms - amount
}

func SwapInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// first time menu or navigated from sub menu
	// show keyboard for the next command
	var msg tgbotapi.Chattable
	if update.CallbackQuery == nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: update.FromChat().ID,
				// TODO: consolidate all telegram send messages
				ReplyMarkup: tokenKeyboard(false),
			},
			Text: "select the source token",
		}
	} else {
		msg = tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			"select the source token", tokenKeyboard(false))
	}

	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-token command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, resp.MessageID, "something went wrong. try again."))
	}
}

func SwapFromTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapInput(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_FROM_TOKEN, fromToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting from-token. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}
	// TODO: move request key expiry to the primary command. its unintuitive to handle it one of the sub-commands. atleast its in the first sub-command!
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving request payload while selecting from-token. %s", err.Error())
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the source network", networkKeyboard(fromToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_FROM_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-network command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func SwapFromNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapFromTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromNetwork := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_FROM_NETWORK, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting from-network. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_TO_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-token command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func SwapToTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapFromNetworkInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_TO_TOKEN, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting to-token. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target network", networkKeyboard(toToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_TO_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-network command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func SwapToNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapToTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_TO_NETWORK, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting to-network. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	// TODO: the keyboard is associated with next sub-command. generalize it for all sub commands
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "How much of your token would you like to swap?", numericKeyboard())
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_QUANTITY,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key tokens command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func SwapQuantiyInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapToNetworkInput(bot, update, false)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "enter") {
		SwapCallback(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_SWAP_CMD_QUANTITY)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching request payload while setting quantity. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_QUANTITY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while setting quantity. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}
	// TODO: handle error
	bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.FromChat().ID, id, "How much of your token would you like to swap?"+quantity,
		numericKeyboard()))

	// the next sub command is still quantity.
	// user completes the command with this subcommand, after pressing "enter"
	// requestKey is no more accessed. will be expired by redis
}

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	SwapInput(bot, update)
}

func SwapCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	requestId := update.CallbackQuery.Message.MessageID
	chatId := update.FromChat().ID
	requestKey := fmt.Sprintf(db.REQUEST_KEY, requestId)

	res := db.RedisClient().HGetAll(context.Background(), requestKey)
	if res.Err() != nil {
		log.Printf("error encountered when fetching swap request payload. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, requestId, "something went wrong. try again."))
		return
	}

	var r SwapRequestInput
	if err := res.Scan(&r); err != nil {
		log.Printf("error parsing swap request payload. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, requestId, "something went wrong. try again."))
		return
	}

	err := swap.SwapTokens(chatId, swap.SwapRequest{
		FromToken:   r.FromToken,
		FromNetwork: r.FromNetwork,
		ToToken:     r.ToToken,
		ToNetwork:   r.ToNetwork,
		Quantity:    r.Quantity,
	})
	if err != nil {
		log.Printf("error executing swap request. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, requestId, "something went wrong. try again."))
		return
	}

	bot.Send(tgbotapi.NewEditMessageText(
		update.FromChat().ID,
		requestId,
		fmt.Sprintf("done! swapped %s tokens from %s:%s to %s:%s",
			r.Quantity, r.FromNetwork, r.FromToken, r.ToNetwork, r.ToToken),
	))
}
