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
const BASE_URL = "https://apigw.okto.tech"

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
	req.Header.Add("x-api-key", os.Getenv("GOOGLE_CLIENT_ID"))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making http req " + err.Error())
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("okto authenticaiton http req not OK")
		return "", errors.New("okto authenticaiton http req not OK")
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto auth response body " + err.Error())
		return "", err
	}

	var authRes AuthResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return "", err
	}

	if authRes.Status != "success" {
		log.Println("okto authenticaiton failed")
		return "", errors.New("okto authenticaiton failed")
	}

	if authRes.Data.Code != 200 {
		log.Println("okto authenticaiton data code not OK")
		return "", errors.New("okto authenticaiton data code not OK")
	}

	return authRes.Data.Token, nil
}