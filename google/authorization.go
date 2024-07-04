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
	req, err := http.NewRequest(http.MethodPost, GOOGLE_DEVICE_CODE_URL, strings.NewReader(fmt.Sprintf("client_id=%s&scope=email%%20profile", os.Getenv("GOOGLE_CLIENT_ID"))))
	if err != nil {
		log.Println("error creating google auth request")
		return DeviceCode{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		log.Printf("error requesting google token id. status code: %d. status: %s", res.StatusCode, res.Status)
		return DeviceCode{}, nil
	}

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading google token response body" + err.Error())
		return DeviceCode{}, err
	}

	var deviceCode DeviceCode
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&deviceCode)
	if err != nil {
		log.Println("error decoding google token response body" + err.Error())
		return DeviceCode{}, err
	}

	return deviceCode, nil
}

func PollAuthorization() (string, error) {
	return "", nil
}
