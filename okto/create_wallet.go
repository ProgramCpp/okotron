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

type Wallet struct {
	NetworkName string `json:"network_name"`
	Address     string `json:"address"`
	Success     bool   `json:"success"`
}

type WalletsData struct {
	Wallets []Wallet `json:"wallets"`
}

type CreateWalletResponse struct {
	Status string      `json:"status"`
	Data   WalletsData `json:"data"`
}

func CreateWallet(authToken string) ([]Wallet, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet", nil)
	if err != nil {
		log.Println("error creating okto create wallet req " + err.Error())
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto create wallet http req " + err.Error())
		return nil, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto create wallet response body " + err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto create wallet http req not OK. " + string(resBytes))
		return nil, errors.New("okto create wallet http req not OK")
	}

	var createWalletRes CreateWalletResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&createWalletRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return nil, err
	}

	if createWalletRes.Status != "success" {
		log.Println("okto request to create wallet failed. " + string(resBytes))
		return nil, errors.New("okto request failed")
	}

	return createWalletRes.Data.Wallets, nil
}
