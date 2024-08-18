package service

import (
	"strings"

	"github.com/pkg/errors"

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
func SwapTokens(r SwapRequest, authToken string, addr string) error {
	transactionPayload, err := lifi.GetQuote(lifi.QuoteRequest{
		FromChain:   okto.NETWORK_NAME_TO_CHAIN_ID[r.FromNetwork],
		FromToken:   r.FromToken,
		ToChain:     okto.NETWORK_NAME_TO_CHAIN_ID[r.ToNetwork],
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
