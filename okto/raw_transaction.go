package okto

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
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

var (
	TXN_IN_PROGRESS = errors.New("txn in progress")
	TXN_FAILED      = errors.New("txn failed")
)

type Order struct {
	OrderId         string `json:"order_id"`
	NetworkName     string `json:"network_name"`
	Status          string `json:"status"`
	TransactionHash string `json:"transaction_hash"`
}

type RawTxnStatusData struct {
	Jobs []Order `json:"jobs"`
}

// TODO: the portfolio response structure is different from whats docuemnted. this is what the api returns. whath out for breaking changes
type RawTxnStatusResponse struct {
	Status string           `json:"status"`
	Data   RawTxnStatusData `json:"data"`
}

func RawTxnStatus(authToken string, jobId string) error {
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/api/v1/rawtransaction/status", nil)
	if err != nil {
		return errors.Wrap(err, "error creating okto transaction status req")
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	params := req.URL.Query()
	params.Add("order_id", jobId)
	req.URL.RawQuery = params.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto transaction status http req " + err.Error())
		return errors.Wrap(err, "")
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto transaction status response body " + err.Error())
		return errors.Wrap(err, "")
	}

	if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto portfolio http req not OK. " + string(resBytes))
		return errors.New("okto transaction status http req not OK")
	}

	var rawTxnStatusRes RawTxnStatusResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&rawTxnStatusRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return errors.Wrap(err, "")
	}

	if rawTxnStatusRes.Status != "success" {
		log.Println("okto request to fetch transaction status failed. " + string(resBytes))
		return errors.New("okto request failed")
	}

	// has only one job
	for _, job := range rawTxnStatusRes.Data.Jobs {
		if job.Status == "SUCCESS" {
			return nil
		} else if job.Status == "" {
			return TXN_IN_PROGRESS
		} else {
			return TXN_FAILED
		}
	}

	return nil
}
