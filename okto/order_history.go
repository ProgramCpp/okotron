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
)

type Job struct {
	OrderId         string `json:"order_id"`
	OrderType       string `json:"order_type"`
	NetworkName     string `json:"network_name"`
	Status          string `json:"status"`
	TransactionHash string `json:"transaction_hash"`
}

func (t Job) String() string {
	// denominate amount in inr, as returned by okto. TODO: convert it to usd?
	return fmt.Sprintf("order id: %s. Network: %s. status: %s. hash: %s",
		t.OrderId, t.NetworkName, t.Status, t.TransactionHash)
}

type OrderHistoryData struct {
	Total int     `json:"total"`
	Jobs  []Job `json:"jobs"`
}

type OrderHistoryResponse struct {
	Status string           `json:"status"`
	Data   OrderHistoryData `json:"data"`
}

// TODO: implement pagination
func OrderHistory(authToken string) ([]Job, error) {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/api/v1/orders", nil)
	if err != nil {
		log.Println("error creating okto order history req " + err.Error())
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto order history http req " + err.Error())
		return nil, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto order history response body " + err.Error())
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		log.Println("okto history http req unauthorized. " + string(resBytes))
		return nil, ERR_UNAUTHORIZED
	} else if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto order history http req not OK. " + string(resBytes))
		return nil, errors.New("okto order history http req not OK")
	}

	var resp OrderHistoryResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&resp)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return nil, err
	}

	if resp.Status != "success" {
		log.Println("okto request to fetch order history failed. " + string(resBytes))
		return nil, errors.New("okto request failed")
	}

	return resp.Data.Jobs, nil
}
