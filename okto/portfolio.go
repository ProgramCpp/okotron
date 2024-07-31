package okto

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

type token struct {
	TokenName    string `json:"token_name"`
	Quantity     string `json:"quantity"`
	AmountInInr  string `json:"amount_in_inr"`
	TokenImage   string `json:"token_image"`
	TokenAddress string `json:"token_address"`
	NetworkName  string `json:"network_name"`
}

type PortfolioData struct {
	Total  string  `json:"total"`
	Tokens []token `json:"tokens"`
}

type portfolioReasponse struct {
	Status string        `json:"status"`
	Data   PortfolioData `json:"data"`
}

func Portfolio(authToken string) ([]token, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/portfolio'", nil)
	if err != nil {
		log.Println("error creating okto portfolio req " + err.Error())
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto portfolio http req " + err.Error())
		return nil, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto portfolio response body " + err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto portfolio http req not OK. " + string(resBytes))
		return nil, errors.New("okto portfolio http req not OK")
	}

	var portfolioRes portfolioReasponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&portfolioRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return nil, err
	}

	if portfolioRes.Status != "success" {
		log.Println("okto request to fetch portfolio failed. " + string(resBytes))
		return nil, errors.New("okto request failed")
	}

	return portfolioRes.Data.Tokens, nil
}
