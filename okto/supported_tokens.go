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

type Token struct {
	TokenName    string `json:"token_name"`
	TokenAddress string `json:"token_address"`
	NetworkName  string `json:"network_name"`
}

func (t Token) String() string {
	return fmt.Sprintf("Token: %s. Address: %s. Network: %s", t.TokenName, t.TokenAddress, t.NetworkName)
}

type TokensData struct {
	Tokens []Token `json:"tokens"`
}

type SupportedTokensResponse struct {
	Status string     `json:"status"`
	Data   TokensData `json:"data"`
}

// osmosis and solana has wallet issues. all the supported tokens are not really supported
func SupportedTokens(authToken string) ([]Token, error) {
	// TODO: paginate to get all supported tokens
	req, err := http.NewRequest(http.MethodGet, BASE_URL+"/api/v1/supported/tokens?page=1&size=10'", nil)
	if err != nil {
		log.Println("error creating okto supported tokens req " + err.Error())
		return nil, err
	}

	req.Header.Add("x-api-key", os.Getenv("OKTO_CLIENT_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error making okto supported tokens http req " + err.Error())
		return nil, err
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading okto supported tokens response body " + err.Error())
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		log.Println("okto supported tokens http req unauthorized. " + string(resBytes))
		return nil, ERR_UNAUTHORIZED
	} else if res.StatusCode != http.StatusOK {
		// TODO: parse error response
		log.Println("okto supported tokens http req not OK. " + string(resBytes))
		return nil, errors.New("okto supported tokens http req not OK")
	}

	var supportedTokensRes SupportedTokensResponse
	err = json.NewDecoder(bytes.NewReader(resBytes)).Decode(&supportedTokensRes)
	if err != nil {
		log.Println("error decoding okto response  " + err.Error())
		return nil, err
	}

	if supportedTokensRes.Status != "success" {
		log.Println("okto request to fetch supported tokens failed. " + string(resBytes))
		return nil, errors.New("okto request failed")
	}

	return supportedTokensRes.Data.Tokens, nil
}