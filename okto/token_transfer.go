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

type TokenTransferRequest struct {
	NetworkName      string
	TokenAddress     string
	Quantity         string
	RecipientAddress string
}

type TransferData struct {
	OrderId string `json:"orderId"`
}

type TransferDataResponse struct {
	Status string       `json:"status"`
	Data   TransferData `json:"data"`
}

func (r TokenTransferRequest) ToReader() (*bytes.Reader, error) {
	buf := bytes.Buffer{}
	e := json.NewEncoder(&buf).Encode(r)
	return bytes.NewReader(buf.Bytes()), e
}

func TokenTransfer(authToken string, r TokenTransferRequest) (string, error) {
	bodyBytes, err := r.ToReader()
	if err != nil {
		log.Println("error serializing okto transfer req " + err.Error())
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/transfers/tokens/execute", bodyBytes)
	if err != nil {
		log.Println("error creating okto transfer req " + err.Error())
		return "", err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "*/*")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto transfer http req " + err.Error())
		return "", err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto transfer response body " + err.Error())
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto transfer http req not OK. " + string(resBytes))
		return "", errors.New("okto transfer http req not OK")
	}

	txfrRes := TransferDataResponse{}
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&txfrRes)
	if err != nil {
		log.Println("error decoding okto transfer response" + err.Error())
		return "", err
	}

	if txfrRes.Status != "success" {
		log.Println("okto transfer request failed. " + string(resBytes))
		// TODO: extract this error string
		return "", errors.New("okto request failed")
	}

	return txfrRes.Data.OrderId, nil
}
