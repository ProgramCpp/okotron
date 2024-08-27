package swap

import (
	"fmt"
	"log"
	"math"
	"strconv"
	
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/programcpp/okotron/okto"
	"github.com/programcpp/okotron/okto/lifi"
)

type SwapRequest struct {
	FromToken   string
	FromNetwork string
	ToToken     string
	ToNetwork   string
	Quantity    string
}

// TODO: okto would support swap directly. there would not be any need for lifi dependability
func SwapTokens(chatId int64, r SwapRequest) error {
	authToken, err := okto.GetAuthToken(chatId)
	if err != nil {
		return errors.Wrap(err, "error fetching okto auth token")
	}

	addr, err := okto.GetAddress(chatId, r.FromNetwork)
	if err != nil {
		return errors.Wrap(err, "error fetching okto address")
	}

	qty, err := strconv.ParseFloat(r.Quantity, 64)
	if err != nil {
		return errors.Wrap(err, "error parsing qty")
	}

	decimal := math.Pow10(okto.TOKEN_TO_DECIMALS[r.FromToken])
	qty *= decimal

	transactionPayload, err := lifi.GetQuote(lifi.QuoteRequest{
		FromChain:   okto.NETWORK_NAME_TO_CHAIN_ID[r.FromNetwork],
		FromToken:   r.FromToken,
		ToChain:     okto.NETWORK_NAME_TO_CHAIN_ID[r.ToNetwork],
		ToToken:     r.ToToken,
		FromAmount:  fmt.Sprintf("%.0f", qty),
		FromAddress: addr,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get a quote for transaction request")
	}
	log.Println(transactionPayload)

	req := okto.RawTxPayload{
		NetworkName: r.FromNetwork,
		Transaction: okto.Transaction{
			From:  gjson.Get(transactionPayload, "from").String(),
			To:    gjson.Get(transactionPayload, "to").String(),
			Data:  gjson.Get(transactionPayload, "data").String(),
			Value: gjson.Get(transactionPayload, "value").String(),
		},
	}
	log.Printf("%v", req)

	_, err = okto.RawTxn(authToken, req)
	if err != nil {
		return errors.Wrap(err, "failed to execute okto transaction")
	}

	return nil
}
