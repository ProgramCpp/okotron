package cmc

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/programcpp/okotron/utils"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

const (
	BASE_URL = "https://pro-api.coinmarketcap.com"
)

// prices are always in INR. for now, only INR is supported
type PricesData struct {
	// map from token to price in INR
	Tokens map[string]float64
}

func NewPricesData() PricesData {
	return PricesData{
		Tokens: make(map[string]float64),
	}
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
func Prices() (PricesData, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Println("error creating cmc quotes http req " + err.Error())
		return PricesData{}, err
	}

	req.Header.Add("X-CMC_PRO_API_KEY", viper.GetString("CMC_KEY"))
	req.Header.Add("Accept", "application/json")

	// TODO: this is hardcoded for now. get supported tokens from cmc as you support more tokens

	tokenSubset := []string{}
	for _, token := range utils.SUPPORTED_TOKENS {
		if !strings.Contains(token, "USD") {
			tokenSubset = append(tokenSubset, token)
		}
	}
	symbols := ""
	for _, token := range tokenSubset {
		symbols += token + ","
	}

	params := url.Values{}
	params.Add("symbol", symbols)
	params.Add("convert", "INR")

	req.URL.RawQuery = params.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making cmc quotes http req " + err.Error())
		return PricesData{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading cmc quotes response body " + err.Error())
		return PricesData{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("cmc quotes http req not OK. " + string(resBytes))
		return PricesData{}, errors.New("cmc quotes http req not OK")
	}

	prices := NewPricesData()
	for _, t := range tokenSubset {
		tokenData := gjson.Get(string(resBytes), "data." + t + ".0") // there is only one currency
		price := gjson.Get(tokenData.String(), "quote.INR.price")
		prices.Tokens[t] = price.Float()
	}

	return prices, nil
}
