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

type RawTxnData struct {
	JobId string `json:"job_id"`
}

// TODO: de-duplicate response struct
type RawTxnResponse struct {
	Status string    `json:"status"`
	Data   RawTxnData `json:"data"`
}

// TODO: de-duplicate http handling
func RawTxn(authToken string, body io.Reader) (RawTxnData, error) {
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/api/v1/rawtransaction/execute", body)
	if err != nil {
		log.Println("error creating okto raw txn req " + err.Error())
		return RawTxnData{}, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "*/*")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto raw txn http req " + err.Error())
		return RawTxnData{}, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto raw txn response body " + err.Error())
		return RawTxnData{}, err
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto raw txn http req not OK. " + string(resBytes))
		return RawTxnData{}, errors.New("okto raw txn http req not OK")
	}

	var rawTxnRes RawTxnResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&rawTxnRes)
	if err != nil {
		log.Println("error decoding okto raw txn response" + err.Error())
		return RawTxnData{}, err
	}

	if rawTxnRes.Status != "success" {
		log.Println("okto request to set raw txn failed. " + string(resBytes))
		// TODO: extract this error string
		return RawTxnData{}, errors.New("okto request failed")
	}
	return rawTxnRes.Data, nil
}
