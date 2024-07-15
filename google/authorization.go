package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	GOOGLE_DEVICE_CODE_URL  = "https://oauth2.googleapis.com/device/code"
	GOOGLE_TOKEN_POLL_URL   = "https://oauth2.googleapis.com/token"
	GOOGLE_DEVICE_SCOPE     = "email%20profile%20openid"
	GOOGLE_OAUTH_GRANT_TYPE = "urn:ietf:params:oauth:grant-type:device_code"
)

var (
	ErrAuthorizationPending = errors.New("authorization_pending")
)

type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	UserCode        string `json:"user_code"`
	VerificationUrl string `json:"verification_url"`
}

func GetDeviceCode() (DeviceCode, error) {
	req, err := http.NewRequest(http.MethodPost, GOOGLE_DEVICE_CODE_URL,
		strings.NewReader(fmt.Sprintf("client_id=%s&scope=%s", os.Getenv("GOOGLE_CLIENT_ID"), GOOGLE_DEVICE_SCOPE)))
	if err != nil {
		log.Println("error creating google auth request")
		return DeviceCode{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error requesting google device code. status code: %d. status: %s", res.StatusCode, res.Status)
		return DeviceCode{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading google device code response body" + err.Error())
		return DeviceCode{}, err
	}

	if res.StatusCode != http.StatusOK {
		var authError AuthError
		err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authError)
		if err != nil {
			log.Println("error decoding google error response" + err.Error())
			return DeviceCode{}, err
		}
		log.Println("google token response not OK. " + authError.ToString())
		return DeviceCode{}, authError
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
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type AuthError struct {
	Error_           string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e AuthError) Error() string {
	return e.Error_
}

func (e AuthError) ToString() string {
	return e.Error_ + "." + e.ErrorDescription
}

func PollAuthorization(deviceCode string) (AccessToken, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID")) // TODO: inject config
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("device_code", deviceCode)
	data.Set("grant_type", GOOGLE_OAUTH_GRANT_TYPE)

	req, err := http.NewRequest(http.MethodPost, GOOGLE_TOKEN_POLL_URL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("error creating google auth request")
		return AccessToken{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error requesting google token id. status code: %d. status: %s", res.StatusCode, res.Status)
		return AccessToken{}, err
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading google token response body" + err.Error())
		return AccessToken{}, err
	}

	// https://developers.google.com/identity/protocols/oauth2/limited-input-device#step-6:-handle-responses-to-polling-requests
	// handle auth errors
	if res.StatusCode != http.StatusOK {
		var authError AuthError
		err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&authError)
		if err != nil {
			log.Println("error decoding google error response. " + string(resBytes) + ". " + err.Error())
			return AccessToken{}, err
		}
		log.Println("google token response not OK. " + authError.ToString())
		return AccessToken{}, authError
	}

	var accessToken AccessToken
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&accessToken)
	if err != nil {
		log.Println("error decoding google token response body" + err.Error())
		return AccessToken{}, err
	}

	return accessToken, nil
}
