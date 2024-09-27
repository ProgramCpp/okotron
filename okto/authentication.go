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

type AuthResponse struct {
	Status string    `json:"status"`
	Data   AuthToken `json:"data"`
}

func (a AuthToken) MarshalBinary() (data []byte, err error) {
	buf := bytes.Buffer{}
	e := json.NewEncoder(&buf).Encode(a)
	return buf.Bytes(), e
}

func Authenticate(idToken string) (AuthToken, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v2/authenticate", strings.NewReader(fmt.Sprintf(
		`
		{
		 	"id_token": "%s"
	 	}
		`, idToken)))
	if err != nil {
		log.Println("error creating okto auth req " + err.Error())
		return AuthToken{}, err
	}

	// TODO: init google client id
	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making http req " + err.Error())
		return AuthToken{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto auth response body " + err.Error())
		return AuthToken{}, err
	}

	fmt.Println(string(resBytes))

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto authenticaiton http req not OK. " + string(resBytes))
		return AuthToken{}, errors.New("okto authenticaiton http req not OK")
	}

	var authRes AuthResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return AuthToken{}, err
	}

	if authRes.Status != "success" {
		log.Println("okto authenticaiton failed")
		return AuthToken{}, errors.New("okto authenticaiton failed. " + string(resBytes))
	}

	return authRes.Data, nil
}
