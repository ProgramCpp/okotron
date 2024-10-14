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

type WalletResponse struct {
	Status string      `json:"status"`
	Data   WalletsData `json:"data"`
}

func (w Wallet) MarshalBinary() (data []byte, err error) {
	buf := bytes.Buffer{}
	e := json.NewEncoder(&buf).Encode(w)
	return buf.Bytes(), e
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

	if res.StatusCode == http.StatusUnauthorized {
		log.Println("okto create wallet http req unauthorized. " + string(resBytes))
		return nil, ERR_UNAUTHORIZED
	} else if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto create wallet http req not OK. " + string(resBytes))
		return nil, errors.New("okto create wallet http req not OK")
	}

	var createWalletRes WalletResponse
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

func Wallets(authToken string) ([]Wallet, error) {	
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/api/v1/wallet'", nil)
	if err != nil {
		log.Println("error creating okto get wallets req " + err.Error())
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto get wallets http req " + err.Error())
		return nil, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto get wallets response body " + err.Error())
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		log.Println("okto get wallets http req unauthorized. " + string(resBytes))
		return nil, ERR_UNAUTHORIZED
	} else if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto get wallets http req not OK. " + string(resBytes))
		return nil, errors.New("okto get wallets http req not OK")
	}

	var WalletsRes WalletResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&WalletsRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return nil, err
	}

	if WalletsRes.Status != "success" {
		log.Println("okto request to fetch get wallets failed. " + string(resBytes))
		return nil, errors.New("okto request failed")
	}

	return WalletsRes.Data.Wallets, nil
}