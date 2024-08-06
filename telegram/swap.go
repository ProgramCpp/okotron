package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/db"
)

const (
	CMD_SWAP_FROM_TOKEN_KEY   = "swap/from-token"
	CMD_SWAP_FROM_NETWORK_KEY = "swap/from-network"
	CMD_SWAP_TO_TOKEN_KEY     = "swap/to-token"
	CMD_SWAP_TO_NETWORK_KEY   = "swap/to-network"
	CMD_SWAP_TO_QUANTITY_KEY  = "swap/quantity"
)

var (
	// TODO: use okto /aupported_tokens and /supported_networks api's
	// do not hardcode networks and tokens
	// for now, all networks returnd by /supported_networks do not work. ex: solana, osmosis
	// an array. do not handle each network separately. do not use enum to treat as first class attributes. oktron is network agnostic
	SUPPORTED_TOKENS   = []string{"APT", "ETH", "MATIC", "USDC", "USDT"}
	SUPPORTED_NETWORKS = map[string][]string{
		"APT":   {"APTOS"},
		"ETH":   {"BASE"},
		"MATIC": {"POLYGON"},
		"USDC":  {"BASE", "POLYGON"},
		"USDT":  {"POLYGON", "APTOS"},
	}
)

type SwapRequest struct {
	FromNetwork string `json:"from_network"`
	ToNetwork   string `json:"to_network"`
	FromToken   string `json:"from_token"`
	ToToken     string `json:"to_token"`
}

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: separate UI componets from backend
	keyboardButtons := []tgbotapi.InlineKeyboardButton{}

	// TODO: see user portfolio for a customized user experience
	for _, n := range SUPPORTED_TOKENS {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var tokenKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	// first time swap menu or navigated from sub menu
	var msg tgbotapi.Chattable
	if update.CallbackQuery == nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      update.FromChat().ID,
				ReplyMarkup: tokenKeyboard,
			},
			Text: "select the source token",
		}
	} else {
		msg = tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID,
			"select the source token", tokenKeyboard)
	}

	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err := db.RedisClient().Set(context.Background(), messageKey, CMD_SWAP_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key from-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func SwapFromToken(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		Swap(bot, update)
		return
	}

	fromToken := update.CallbackQuery.Data
	id := update.CallbackQuery.Message.MessageID
	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_FROM_TOKEN_KEY, fromToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
	}

	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("back", "back"),
	}

	for _, n := range SUPPORTED_NETWORKS[fromToken] {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the source network", networkKeyboard)
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_SWAP_CMD_FROM_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key from-network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func SwapFromNetwork(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapFromToken(bot, update, false)
		return
	}

	fromNetwork := update.CallbackQuery.Data
	id := update.CallbackQuery.Message.MessageID

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_FROM_NETWORK_KEY, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
	}

	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("back", "back"),
	}

	for _, n := range SUPPORTED_TOKENS {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var tokenKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard)
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_SWAP_CMD_TO_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key to-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func SwapToToken(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapFromNetwork(bot, update, false)
		return
	}

	toToken := update.CallbackQuery.Data
	id := update.CallbackQuery.Message.MessageID

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_TOKEN_KEY, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting to-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
	}

	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("back", "back"),
	}

	for _, n := range SUPPORTED_NETWORKS[toToken] {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target network", networkKeyboard)
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_SWAP_CMD_TO_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key to-network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func SwapToNetwork(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapToToken(bot, update, false)
		return
	}

	toToken := update.CallbackQuery.Data
	id := update.CallbackQuery.Message.MessageID

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_NETWORK_KEY, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting to-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
	}

	var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("7", "7"),
			tgbotapi.NewInlineKeyboardButtonData("8", "8"),
			tgbotapi.NewInlineKeyboardButtonData("9", "9"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("back", "back"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:", networkKeyboard)
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_SWAP_CMD_TOKENS,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key tokens command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func SwapQuantiy(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		SwapToNetwork(bot, update, false)
		return
	}

	// user has input all the request params. process the request
	if update.CallbackQuery.Data == "enter" {
		swapTokens()
	}

	quantity := update.CallbackQuery.Data
	id := update.CallbackQuery.Message.MessageID

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_QUANTITY_KEY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while setting quantity. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
	}

	var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("7", "7"),
			tgbotapi.NewInlineKeyboardButtonData("8", "8"),
			tgbotapi.NewInlineKeyboardButtonData("9", "9"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("back", "back"),
			tgbotapi.NewInlineKeyboardButtonData("enter", "enter"),
		),
	)

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:"+quantity, networkKeyboard)
	// TODO: handle error
	bot.Send(msg)

	// the next sub command is still quantity. wait until done
}


func swapTokens(){
	
}