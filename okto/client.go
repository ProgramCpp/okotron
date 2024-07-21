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
	"strings"
)

// TODO: make this configurable
// Sandbox : https://sandbox-api.okto.tech
// Staging : https://3p-bff.oktostage.com
// Production : https://apigw.okto.tech
const BASE_URL = "https://sandbox-api.okto.tech" // TODO: make it configurable for different environments

type AuthData struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Action  string `json:"action"`
	Code    int    `json:"code"`
}

type AuthResponse struct {
	Status string   `json:"status"`
	Data   AuthData `json:"data"`
}

func Authenticate(idToken string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/authenticate", strings.NewReader(fmt.Sprintf(
		`
		{
		 	"id_token": "%s"
	 	}
		`, idToken)))
	if err != nil {
		log.Println("error creating okto auth req " + err.Error())
		return "", err
	}

	// TODO: init google client id
	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making http req " + err.Error())
		return "", err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto auth response body " + err.Error())
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto authenticaiton http req not OK. " + string(resBytes))
		return "", errors.New("okto authenticaiton http req not OK")
	}

	var authRes AuthResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return "", err
	}

	if authRes.Status != "success" {
		log.Println("okto authenticaiton failed")
		return "", errors.New("okto authenticaiton failed. " + string(resBytes))
	}

	if authRes.Data.Code != 200 {
		log.Println("okto authenticaiton data code not OK")
		return "", errors.New("okto authenticaiton data code not OK")
	}

	return authRes.Data.Token, nil
}

type Token struct {
	TokenName    string `json:"token_name"`
	TokenAddress string `json:"token_address"`
	NetworkName  string `json:"network_name"`
}

func (t Token) String() string {
	return fmt.Sprintf("Token: %s. Address: %s. Network: %s", t.TokenName, t.TokenAddress, t.NetworkName)
}

type TokensData struct {
	Tokens []Token `json:"tokens"`
}

type SupportedTokensResponse struct {
	Status string     `json:"status"`
	Data   TokensData `json:"data"`
}

func SupportedTokens(authToken string) ([]Token, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/api/v1/supported/tokens?page=1&size=10'", nil)
	if err != nil {
		log.Println("error creating okto supported tokens req " + err.Error())
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto supported tokens http req " + err.Error())
		return nil, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto supported tokens response body " + err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto authenticaiton http req not OK. " + string(resBytes))
		return nil, errors.New("okto authenticaiton http req not OK")
	}

	var supportedTokensRes SupportedTokensResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&supportedTokensRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return nil, err
	}

	if supportedTokensRes.Status != "success" {
		log.Println("okto request to fetch supported tokens failed. " + string(resBytes))
		return nil, errors.New("okto request failed")
	}

	return supportedTokensRes.Data.Tokens, nil
}
