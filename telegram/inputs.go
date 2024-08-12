package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
)

func TokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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
		log.Printf("error encountered when saving message key from-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func FromToken(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		TokenInput(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_ANY_CMD_FROM_TOKEN, fromToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting from-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
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

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_ANY_CMD_FROM_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func FromNetwork(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		FromToken(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromNetwork := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_ANY_CMD_FROM_NETWORK, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard())
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_ANY_CMD_TO_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func ToToken(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		FromNetwork(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_ANY_CMD_TO_TOKEN, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting to-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target network", networkKeyboard(toToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_ANY_CMD_TO_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func ToNetwork(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		ToToken(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_ANY_CMD_TO_NETWORK, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting to-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	// TODO: the keyboard is associated with next sub-command. generalize it for all sub commands
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:", numericKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_ANY_CMD_QUANTITY,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key tokens command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func Quantiy(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		ToNetwork(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	// user has input all the request params. process the request
	if update.CallbackQuery.Data == "enter" {
		Swap(bot, update)
		return
	}

	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_ANY_CMD_QUANTITY)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching request payload while setting quantity. %s", res.Err())
		Send(bot, update, "something went wrong. try again.")
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_ANY_CMD_QUANTITY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while setting quantity. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	// TODO: handle error
	bot.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.FromChat().ID, id, "enter token quantity:"+quantity,
		numericKeyboard(true)))

	// the next sub command is still quantity.
	// user completes the command with this subcommand, after pressing "enter"
	// requestKey is no more accessed. will be expired by redis
}