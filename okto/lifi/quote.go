package lifi

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

const BASE_URL = "https://li.quest"

type QuoteRequest struct {
	FromChain   string
	ToChain     string
	FromToken   string
	ToToken     string
	FromAmount  string
	FromAddress string
}

func GetQuote(r QuoteRequest) (string, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/v1/quote", nil)
	if err != nil {
		log.Println("error creating lifi get-quote http req " + err.Error())
		return "", err
	}

	params := url.Values{}
	params.Add("fromChain", r.FromChain)
	params.Add("fromToken", r.FromToken)
	params.Add("toChain", r.ToChain)
	params.Add("toToken", r.ToToken)
	params.Add("fromAddress", r.FromAddress)
	params.Add("fromAmount", r.FromAmount)

	req.URL.RawQuery = params.Encode()
	// req.Header.Add("Content-Type", "application/json") // not accepted

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making lifi get-quote http req " + err.Error())
		return "", err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading lifi get-quote response body " + err.Error())
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("lifi get-quote http req not OK. " + string(resBytes))
		return "", errors.New("lifi get-quote http req not OK")
	}

	result := gjson.Get(string(resBytes), "transactionRequest")

	return result.String(), nil
}
