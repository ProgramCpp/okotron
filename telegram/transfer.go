package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// redis tags are defined as constants. keep in SYNC. the keys are saved with COMMAND Names and deserialized with redis tags!
type TransferRequestInput struct {
	FromToken   string `json:"from_token" redis:"/transfer/from-token"`
	FromNetwork string `json:"from_network" redis:"/transfer/from-network"`
	Quantity    string `json:"quantity" redis:"/transfer/quantity"`
	Address     string `json:"address" redis:"/transfer/address"`
}

func TransferInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: handle error
	resp, _ := bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: update.FromChat().ID,
			ReplyMarkup: tgbotapi.ForceReply{
				ForceReply: true,
			},
		},
		Text: "enter to-address:",
	})

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_TRANSFER_CMD_ADDRESS,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-token command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, resp.MessageID, "something went wrong. try again."))
	}
}

func TransferAddressInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	id := update.Message.ReplyToMessage.MessageID
	address := update.Message.Text

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

	id = resp.MessageID

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, id)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_TRANSFER_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-token command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, resp.MessageID, "something went wrong. try again."))
	}

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err = db.RedisClient().HSet(context.Background(), requestKey, CMD_TRANSFER_CMD_ADDRESS, address).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while saving address. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
		return
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
		log.Printf("error encountered when saving transfer message key tokens command. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, id, "something went wrong. try again."))
	}
}

func TransferQuantityInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		TransferFromTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	if strings.Contains(update.CallbackQuery.Data, "enter") {
		TransferCallback(bot, update)
		return
	}

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

func TransferCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.CallbackQuery.Message.MessageID
	chatId := update.FromChat().ID
	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)

	res := db.RedisClient().HGetAll(context.Background(), requestKey)
	if res.Err() != nil {
		log.Printf("error encountered when fetching transfer request payload. %s", res.Err())
		bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "something went wrong. try again."))
		return
	}

	var r TransferRequestInput
	if err := res.Scan(&r); err != nil {
		log.Printf("error parsing transfer request payload. %s", err.Error())
		bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "something went wrong. try again."))
		return
	}

	tokAddr := okto.TOKEN_TO_NETWORK_TO_ADDRESS[r.FromToken][r.FromNetwork]
	if tokAddr == okto.NATIVE_TOKEN_ADDR {
		tokAddr = " " // empty string for okto transfer
	}

	_, err := okto.TokenTransfer(chatId, okto.TokenTransferRequest{
		NetworkName:      r.FromNetwork,
		TokenAddress:     tokAddr,
		Quantity:         r.Quantity,
		RecipientAddress: r.Address,
	})
	if err != nil && errors.Is(err, okto.ERR_UNAUTHORIZED) {
		log.Printf("error executing transfer request. " + err.Error())
		Send(bot, update, "unauthorized. login and try again.")
		return
	} else if err != nil {
		log.Printf("error executing transfer request. %s", err.Error())
		bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "something went wrong. try again."))
		return
	}

	bot.Send(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatId,
			MessageID:   id,
			ReplyMarkup: nil,
		},
		Text: fmt.Sprintf("done! transfer %s tokens from %s:%s to %s",
			r.Quantity, r.FromNetwork, r.FromToken, r.Address),
	})
}

func Transfer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	TransferInput(bot, update)
}
