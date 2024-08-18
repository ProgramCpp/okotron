package cmc

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/programcpp/okotron/utils"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

const (
	BASE_URL = " https://pro-api.coinmarketcap.com"
)

// prices are always in INR. for now, only INR is supported
type Price struct {
	// map from token to price in INR
	Tokens map[string]float64
}

/*
references:
CMC fiat currency ids:
https://coinmarketcap.com/api/documentation/v1/#section/Standards-and-Conventions
CMC token symbols:
https://stackoverflow.com/a/70028568/2508038
CMC quote api:
https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyQuotesLatest
*/
func Prices() (Price, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Println("error creating cmc quotes http req " + err.Error())
		return Price{}, err
	}

	// TODO: this is hardcoded for now. get supported tokens from cmc as you support more tokens
	symbols := ""
	for _, token := range utils.SUPPORTED_TOKENS {
		symbols += token + ","
	}

	params := url.Values{}
	params.Add("symbol", symbols)
	params.Add("convert", "2796") // this is INR. for now, only INR is supported
	params.Add("X-CMC_PRO_API_KEY", viper.GetString("CMC_KEY"))
	params.Add("Accept", "application/json")

	req.URL.RawQuery = params.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making cmc quotes http req " + err.Error())
		return Price{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading cmc quotes response body " + err.Error())
		return Price{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("cmc quotes http req not OK. " + string(resBytes))
		return Price{}, errors.New("cmc quotes http req not OK")
	}

	tokens := gjson.GetMany(string(resBytes), "data.*")

	prices := Price{}
	for _, t := range tokens {
		symbol := gjson.Get(t.String(), "symbol")
		price := gjson.Get(t.String(), "quote.INR.price")
		prices.Tokens[symbol.String()] = price.Float()
	}

	return prices, nil
}
