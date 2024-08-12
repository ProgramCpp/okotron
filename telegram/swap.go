package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
	"github.com/programcpp/okotron/okto/lifi"

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

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update){
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
