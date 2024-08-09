package okto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type PortfolioTokenInfo struct {
	TokenName    string `json:"token_name"`
	Quantity     string `json:"quantity"`
	AmountInInr  string `json:"amount_in_inr"`
	TokenImage   string `json:"token_image"`
	TokenAddress string `json:"token_address"`
	NetworkName  string `json:"network_name"`
}

func (t PortfolioTokenInfo) String() string {
	return fmt.Sprintf("token: %s. Network: %s. Quantity: %s. amount: %s.",
		t.TokenName, t.NetworkName, t.Quantity, t.AmountInInr)
}

type PortfolioData struct {
	Total  int                  `json:"total"`
	Tokens []PortfolioTokenInfo `json:"tokens"`
}

// TODO: the portfolio response structure is different from whats docuemnted. this is what the api returns. whath out for breaking changes
type portfolioReasponse struct {
	Status string        `json:"status"`
	Data   PortfolioData `json:"data"`
}

func Portfolio(authToken string) ([]PortfolioTokenInfo, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/api/v1/portfolio", nil)
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
