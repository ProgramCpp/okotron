package swap

import (
	"fmt"
	"log"
	"math"
	"math/big"
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

	fromTokAddr := okto.TOKEN_TO_NETWORK_TO_ADDRESS[r.FromToken][r.FromNetwork]

	transactionPayload, err := lifi.GetQuote(lifi.QuoteRequest{
		FromChain:   okto.NETWORK_NAME_TO_CHAIN_ID[r.FromNetwork],
		FromToken:   fromTokAddr,
		ToChain:     okto.NETWORK_NAME_TO_CHAIN_ID[r.ToNetwork],
		ToToken:     okto.TOKEN_TO_NETWORK_TO_ADDRESS[r.ToToken][r.ToNetwork],
		FromAmount:  fmt.Sprintf("%.0f", qty),
		FromAddress: addr,
	})
	if err != nil {
		return errors.Wrap(err, "failed to get a quote for transaction request")
	}
	log.Println(transactionPayload)

	if fromTokAddr != okto.NATIVE_TOKEN_ADDR {
		err = okto.ApproveTokenTransfer(
			authToken, 
			r.FromNetwork, 
			fromTokAddr, 
			gjson.Get(transactionPayload, "to").String(), 
			big.NewInt(int64(qty)), 
			addr)
		if err != nil {
			return errors.Wrap(err, "failed to approve transaction")
		}
	}

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
