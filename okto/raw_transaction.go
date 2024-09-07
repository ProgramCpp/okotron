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

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Data  string `json:"data"`
	Value string `json:"value"`
}
type RawTxPayload struct {
	NetworkName string      `json:"network_name"`
	Transaction Transaction `json:"transaction"`
}

type RawTxnResponseData struct {
	JobId string `json:"job_id"`
}

// TODO: de-duplicate response struct
type RawTxnResponse struct {
	Status string             `json:"status"`
	Data   RawTxnResponseData `json:"data"`
}

// TODO: de-duplicate http handling
func RawTxn(authToken string, r RawTxPayload) (RawTxnResponseData, error) {
	bodyBytes := bytes.Buffer{}
	err := json.NewEncoder(&bodyBytes).Encode(r)
	if err != nil {
		log.Println("error encoding transaction payload " + err.Error())
		return RawTxnResponseData{}, err
	}

	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/rawtransaction/execute", &bodyBytes)
	if err != nil {
		log.Println("error creating okto raw txn req " + err.Error())
		return RawTxnResponseData{}, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "*/*")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto raw txn http req " + err.Error())
		return RawTxnResponseData{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto raw txn response body " + err.Error())
		return RawTxnResponseData{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto raw txn http req not OK. " + string(resBytes))
		return RawTxnResponseData{}, errors.New("okto raw txn http req not OK")
	}

	var rawTxnRes RawTxnResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&rawTxnRes)
	if err != nil {
		log.Println("error decoding okto raw txn response" + err.Error())
		return RawTxnResponseData{}, err
	}

	if rawTxnRes.Status != "success" {
		log.Println("okto request to execute raw txn failed. " + string(resBytes))
		// TODO: extract this error string
		return RawTxnResponseData{}, errors.New("okto request failed")
	}
	return rawTxnRes.Data, nil
}

func RawTxnStatus(authToken string, jobId string) error {
	return nil
}