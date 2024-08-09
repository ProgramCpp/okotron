package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
	"github.com/programcpp/okotron/okto/lifi"
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
type SwapRequestInput struct {
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

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
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

	requestKey := fmt.Sprintf(db.SWAP_REQUEST_KEY, id)
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

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
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

	requestKey := fmt.Sprintf(db.SWAP_REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_FROM_NETWORK_KEY, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard())
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
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

	requestKey := fmt.Sprintf(db.SWAP_REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_TO_TOKEN_KEY, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving swap request payload while selecting to-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target network", networkKeyboard(toToken))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
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

	requestKey := fmt.Sprintf(db.SWAP_REQUEST_KEY, id)
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

	messageKey := fmt.Sprintf(db.MESSAGE_KEY, resp.MessageID)
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

	requestKey := fmt.Sprintf(db.SWAP_REQUEST_KEY, id)
	// user has input all the request params. process the request
	if update.CallbackQuery.Data == "enter" {
		res := db.RedisClient().HGetAll(context.Background(), requestKey)
		if res.Err() != nil {
			log.Printf("error encountered when fetching swap request payload. %s", res.Err())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		var r SwapRequestInput
		if err := res.Scan(&r); err != nil {
			log.Printf("error parsing swap request payload. %s", err.Error())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		tokRes := db.RedisClient().Get(context.Background(), fmt.Sprintf(db.OKTO_AUTH_TOKEN_KEY, update.Message.Chat.ID))
		if tokRes.Err() != nil {
			log.Printf("error fetching okto auth token. %s", tokRes.Err())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		addrRes := db.RedisClient().Get(context.Background(), fmt.Sprintf(db.OKTO_ADDRESSES_KEY, update.Message.Chat.ID))
		if addrRes.Err() != nil {
			log.Printf("error fetching okto auth token. %s", addrRes.Err())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		var wallets []okto.Wallet
		err := json.NewDecoder(strings.NewReader(addrRes.Val())).Decode(wallets)
		if addrRes.Err() != nil {
			log.Printf("error decoding okto wallets. %s", addrRes.Err())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		err = swapTokens(r, tokRes.Val(), addrRes.Val())
		if err != nil {
			log.Printf("error executing swap request. %s", err.Error())
			Send(bot, update, "something went wrong. try again.")
			return
		}

		bot.Send(tgbotapi.NewEditMessageText(
			update.FromChat().ID,
			id,
			fmt.Sprintf("done! swapped %s tokens from %s:%s to %s:%s",
				r.Quantity, r.FromNetwork, r.FromToken, r.ToNetwork, r.ToToken),
		))
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

// TODO: separate telegram and service concerns
func swapTokens(r SwapRequestInput, authToken string, addr string) error {
	transactionPayload, err := lifi.GetQuote(lifi.QuoteRequest{
		FromChain:   r.FromNetwork,
		FromToken:   r.FromToken,
		ToChain:     r.ToNetwork,
		ToToken:     r.ToToken,
		FromAmount:  r.Quantity,
		FromAddress: addr,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get a quote for transaction request")
	}

	_, err = okto.RawTxn(authToken, strings.NewReader(transactionPayload), r.FromNetwork)
	if err != nil {
		return errors.Wrap(err, "failed to execute okto transaction")
	}

	return nil
}
