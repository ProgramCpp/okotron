package telegram

import (
	"encoding/json"
	"strings"

	"github.com/programcpp/okotron/okto"
)

func getAuthToken(authTokenjson string) string {
	// TODO: handle token not found
	authToken := okto.AuthToken{}
	// TODO: handle decoding error
	json.NewDecoder(strings.NewReader(authTokenjson)).Decode(&authToken)
	return authToken.AuthToken
}