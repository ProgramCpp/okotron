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

func RefreshTokens(authToken AuthToken) (AuthToken, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/refresh_token", nil)
	if err != nil {
		log.Println("error creating okto refresh token req " + err.Error())
		return AuthToken{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "*/*")
	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Authorization", "Bearer "+authToken.AuthToken)
	req.Header.Add("x-refresh-authorization", authToken.RefreshAuthToken)
	req.Header.Add("x-device-token", authToken.DeviceToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto refresh token req " + err.Error())
		return AuthToken{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto refresh token body " + err.Error())
		return AuthToken{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto refresh token http req not OK. " + string(resBytes))
		return AuthToken{}, errors.New("okto refresh token http req not OK")
	}

	var authTokenRes AuthTokenResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authTokenRes)
	if err != nil {
		log.Println("error decoding okto refresh token response" + err.Error())
		return AuthToken{}, err
	}

	if authTokenRes.Status != "success" {
		log.Println("okto request to refresh token failed. " + string(resBytes))
		// TODO: remove the expired auth token, so that user can authorize okotron again
		return AuthToken{}, errors.New("okto request failed")
	}
	return authTokenRes.Data, nil
}
