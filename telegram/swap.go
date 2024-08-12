package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
	"github.com/programcpp/okotron/okto/lifi"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// redis tags are defined as constants. keep in SYNC. the keys as saved with COMMAND Names and deserialized with redis tags!
type SwapRequestInput struct {
	FromToken   string `json:"from_token" redis:"any/from-token"`
	FromNetwork string `json:"from_network" redis:"any/from-network"`
	ToToken     string `json:"to_token" redis:"any/to-token"`
	ToNetwork   string `json:"to_network" redis:"any/to-network"`
	Quantity    string `json:"quantity" redis:"any/quantity"` // use okto terms - quantity. not lifi terms - amount
}

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

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_FROM_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func FromTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		TokenInput(bot, update)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_FROM_TOKEN, fromToken).Err()
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

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_FROM_NETWORK,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key from-network command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func FromNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		FromTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	fromNetwork := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_FROM_NETWORK, fromNetwork).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting from-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "select the target token", tokenKeyboard())
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_TO_TOKEN,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving message key to-token command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func ToTokenInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		FromNetworkInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_TO_TOKEN, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting to-token. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
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
		Send(bot, update, "something went wrong. try again.")
	}
}

func ToNetworkInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		ToTokenInput(bot, update, false)
		return
	}

	id := update.CallbackQuery.Message.MessageID
	toToken := update.CallbackQuery.Data

	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_TO_NETWORK, toToken).Err()
	if err != nil {
		log.Printf("error encountered when saving request payload while selecting to-network. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
		return
	}

	// TODO: the keyboard is associated with next sub-command. generalize it for all sub commands
	msg := tgbotapi.NewEditMessageTextAndMarkup(update.FromChat().ID, id, "enter token quantity:", numericKeyboard(true))
	// TODO: handle error
	resp, _ := bot.Send(msg)

	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err = db.RedisClient().Set(context.Background(), subcommandKey, CMD_SWAP_CMD_QUANTITY,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving swap message key tokens command. %s", err.Error())
		Send(bot, update, "something went wrong. try again.")
	}
}

func QuantiyInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	if isBack {
		ToNetworkInput(bot, update, false)
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
		Send(bot, update, "something went wrong. try again.")
		return
	} else if res.Err() != redis.Nil {
		quantity = res.Val() + quantity
	}

	err := db.RedisClient().HSet(context.Background(), requestKey, CMD_SWAP_CMD_QUANTITY, quantity).Err()
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

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	TokenInput(bot, update)
}

func SwapCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.CallbackQuery.Message.MessageID
	requestKey := fmt.Sprintf(db.REQUEST_KEY, id)
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
