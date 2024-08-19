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

func filterStableCoins(tokens []string) []string {
	tokenSet := []string{}
	for _, token := range tokens {
		if !strings.Contains(token, "USD") {
			tokenSet = append(tokenSet, token)
		}
	}

	return tokenSet
}

func formatTokenSymbols(tokens []string) string {
	symbols := ""
	for _, token := range tokens {
		symbols += token + ","
	}
	return symbols
}

func getPrices(req *http.Request) (PricesData, error){
	req.Header.Add("X-CMC_PRO_API_KEY", viper.GetString("CMC_KEY"))
	req.Header.Add("Accept", "application/json")

	// TODO: this is hardcoded for now. get supported tokens from cmc/ okto as you support more tokens
	tokenSet := filterStableCoins(utils.SUPPORTED_TOKENS)

	symbols := formatTokenSymbols(tokenSet)
	params := req.URL.Query()
	params.Add("symbol", symbols)

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
	for _, t := range tokenSet {
		tokenData := gjson.Get(string(resBytes), "data."+t+".0") // there is only one currency
		price := gjson.Get(tokenData.String(), "quote.INR.price")
		prices.Tokens[t] = price.Float()
	}

	return prices, nil
}

func PricesInCurrency() (PricesData, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Println("error creating cmc quotes http req " + err.Error())
		return PricesData{}, err
	}

	params := url.Values{}
	params.Add("convert", "INR")

	req.URL.RawQuery = params.Encode()

	return getPrices(req)
}

func PricesInTokens() (PricesData, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Println("error creating cmc quotes http req " + err.Error())
		return PricesData{}, err
	}

	tokenSet := filterStableCoins(utils.SUPPORTED_TOKENS)
	symbols := formatTokenSymbols(tokenSet)

	params := url.Values{}
	params.Add("convert", symbols)

	req.URL.RawQuery = params.Encode()

	return getPrices(req)
}