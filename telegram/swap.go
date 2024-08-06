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
	// an array. do not handle each network separately. do not use enum to treat as first class attributes. okotron is network agnostic
	SUPPORTED_TOKENS   = []string{"APT", "ETH", "MATIC", "USDC", "USDT"}
	SUPPORTED_NETWORKS = map[string][]string{
		"APT":   {"APTOS"},
		"ETH":   {"BASE"},
		"MATIC": {"POLYGON"},
		"USDC":  {"BASE", "POLYGON"},
		"USDT":  {"POLYGON", "APTOS"},
	}
)

// redis tags are defined above as constants. keep in SYNC
type SwapRequest struct {
	FromToken   string `json:"from_token" redis:"swap/from-token"`
	FromNetwork string `json:"from_network" redis:"swap/from-network"`
	ToToken     string `json:"to_token" redis:"swap/to-token"`
	ToNetwork   string `json:"to_network" redis:"swap/to-network"`
	Quantity    string `json:"quantity" redis:"swap/quantity"`
}

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// first time swap menu or navigated from sub menu
	// show keyboard for the next command
	var msg tgbotapi.Chattable
	if update.CallbackQuery == nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      update.FromChat().ID,
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

	id := update.CallbackQuery.Message.MessageID
	fromToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_FROM_TOKEN_KEY, fromToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}
	// TODO: move this to the primary command -swap. its unintuitive to handle it one of the sub-commands
	_, err = db.RedisClient().Expire(context.Background(), requestKey, time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Result()
	if err != nil {
		// just logging for now. this will result in stale values. do not stop user flow
		log.Printf("error encountered when saving swap request payload while selecting from-token. %s", err.Error())
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the source network", networkKeyboard(fromToken))
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

	id := update.CallbackQuery.Message.MessageID
	fromNetwork := update.CallbackQuery.Data

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_FROM_NETWORK_KEY, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard())
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

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_TOKEN_KEY, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting to-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target network", networkKeyboard(toToken))
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

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf("swap_%d", id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_NETWORK_KEY, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting to-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	// TODO: the keyboard is associated with next sub-command. generalize it for all sub commands
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:", numericKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf("message_%d", resp.MessageID)
	err = db.RedisClient().Set(context.Background(), messageKey, CMD_SWAP_CMD_QUANTITY,
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

	id := update.CallbackQuery.Message.MessageID
	quantity := update.CallbackQuery.Data

	requestKey := fmt.Sprintf("swap_%d", id)
	// user has input all the request params. process the request
	if update.CallbackQuery.Data == "enter" {
		res := db.RedisClient().HGetAll(context.Background(), requestKey)
		if res.Err() != nil {
			log.Printf("error encountered when fetching swap request payload. %s", res.Err())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		var r SwapRequest
		if err := res.Scan(&r); err != nil {
			log.Printf("error parsing swap request payload. %s", err.Error())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		bot.Send(tgbotapi.NewEditMessageText(
			update.FromChat().ID,
			id,
			fmt.Sprintf("done! swapped %s tokens from %s:%s to %s:%s",
				r.Quantity, r.FromNetwork, r.FromToken, r.ToNetwork, r.ToToken),
		))

		swapTokens(r)
		return
	}

	// handle first digit of quantity. redis returns Nil error of the field is not found
	res := db.RedisClient().HGet(context.Background(), requestKey, CMD_SWAP_TO_QUANTITY_KEY)
	if res.Err() != nil && res.Err() != redis.Nil {
		log.Printf("error encountered when fetching swap request payload while setting quantity. %s", res.Err())
		Send(bot, update, "something went wrong. try again.")
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_QUANTITY_KEY, quantity).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while setting quantity. %s", err.Error())
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

func swapTokens(r SwapRequest) {

}
