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

type AuthToken struct {
	AuthToken        string `json:"auth_token"`
	Message          string `json:"message"`
	RefreshAuthToken string `json:"refresh_auth_token"`
	DeviceToken      string `json:"device_token"`
}

type AuthTokenResponse struct {
	Status string    `json:"status"`
	Data   AuthToken `json:"data"`
}

func SetPin(idToken, token, pin string) (AuthToken, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/wallet", strings.NewReader(fmt.Sprintf(
		`
		{
		 	"id_token": "%s",
			"token": "%s",
            "relogin_pin": "%s",
            "purpose": "set_pin"
	 	}
		`, idToken, token, pin)))
	if err != nil {
		log.Println("error creating okto pin req " + err.Error())
		return AuthToken{}, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto pin http req " + err.Error())
		return AuthToken{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto pin response body " + err.Error())
		return AuthToken{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto pin http req not OK. " + string(resBytes))
		return AuthToken{}, errors.New("okto pin http req not OK")
	}

	var authTokenRes AuthTokenResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authTokenRes)
	if err != nil {
		log.Println("error decoding okto pin response" + err.Error())
		return AuthToken{}, err
	}

	if authTokenRes.Status != "success" {
		log.Println("okto request to set pin failed. " + string(resBytes))
		return AuthToken{}, errors.New("okto request failed")
	}
	return authTokenRes.Data, nil
}
