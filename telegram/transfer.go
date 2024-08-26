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

func TransferInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// show keyboard for the next command
	// TODO: handle error
	resp, _ := bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: update.FromChat().ID,
			// TODO: consolidate all telegram send messages
			ReplyMarkup: tokenKeyboard(false),
		},
		Text: "select the token to transfer:",
	})

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_TRANSFER_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-token command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, resp.MessageID, "something went wrong. try again."))
	}
}

func TransferFromTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		TransferInput(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_TRANSFER_CMD_FROM_TOKEN, fromToken).Err()
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

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the network", networkKeyboard(fromToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_TRANSFER_CMD_FROM_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-network command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func TransferFromNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		TransferInput(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromNetwork := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_TRANSFER_CMD_FROM_NETWORK, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting from-network. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:", numericKeyboard())
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_TRANSFER_CMD_QUANTITY,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key tokens command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func TransferQuantityInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		TransferFromTokenInput(bot, update, false)
		return
	}

	if strings.Contains(update.CallbackQuery.Data, "enter") {
		TransferCallback(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_TRANSFER_CMD_QUANTITY)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching request payload while setting quantity. %s", res.Err())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_TRANSFER_CMD_QUANTITY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while setting quantity. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
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

func TransferCallback (bot *tgbotapi.BotAPI, update tgbotapi.Update) {
}

func Transfer (bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	TransferInput(bot, update)
}