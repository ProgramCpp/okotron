package google

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

const (
	GOOGLE_DEVICE_CODE_URL = "https://oauth2.googleapis.com/device/code"
	GOOGLE_TOKEN_POLL_URL  = "https://oauth2.googleapis.com/token"
	GOOGLE_DEVICE_SCOPE    = "email%20profile%20openid"
)

var (
	ErrAuthorizationPending = errors.New("pending user authorization")
)

type DeviceCode struct {
	DeviceCode      string `json:device_code`
	ExpiresIn       int    `json:expires_in`
	Interval        int    `json:interval`
	UserCode        string `json:user_code`
	VerificationUrl string `json:verification_url`
}

func GetDeviceCode() (DeviceCode, error) {
	req, err := http.NewRequest(http.MethodPost, GOOGLE_DEVICE_CODE_URL,
		strings.NewReader(fmt.Sprintf("client_id=%s&scope=%s", os.Getenv("GOOGLE_CLIENT_ID"), GOOGLE_DEVICE_SCOPE)))
	if err != nil {
		log.Println("error creating google auth request")
		return DeviceCode{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	// TODO: 403 returns error in response body
	if err != nil || res.StatusCode != http.StatusOK {
		log.Printf("error requesting google device code. status code: %d. status: %s", res.StatusCode, res.Status)
		return DeviceCode{}, nil
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading google device code response body" + err.Error())
		return DeviceCode{}, err
	}

	var deviceCode DeviceCode
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&deviceCode)
	if err != nil {
		log.Println("error decoding google device code response body" + err.Error())
		return DeviceCode{}, err
	}

	return deviceCode, nil
}

type AccessToken struct {
	AccessToken  string `json:access_token`
	ExpiresIn    int    `json:expires_in`
	IdToken      string `json: id_token`
	Scope        string `json:scope`
	TokenType    string `json:token_type`
	RefreshToken string `json:refresh_token`
}

func PollAuthorization(deviceCode string) (AccessToken, error) {
	reqBody := fmt.Sprintf(
		`
		client_id=%s&
		client_secret=%s&
		device_code=%s&
		grant_type=urn%%3Aietf%%3Aparams%%3Aoauth%%3Agrant-type%%3Adevice_code
		`, os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), deviceCode)
	req, err := http.NewRequest(http.MethodPost, GOOGLE_TOKEN_POLL_URL, strings.NewReader(reqBody))
	if err != nil {
		log.Println("error creating google auth request")
		return AccessToken{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Printf("error requesting google token id. status code: %d. status: %s", res.StatusCode, res.Status)
		return AccessToken{}, nil
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading google token response body" + err.Error())
		return AccessToken{}, err
	}

	var accessToken AccessToken
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&accessToken)
	if err != nil {
		log.Println("error decoding google token response body" + err.Error())
		return AccessToken{}, err
	}

	return accessToken, nil
}
