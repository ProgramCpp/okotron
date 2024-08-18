package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
)

func GetAuthToken(chatId int64) (string, error) {
	authTokenKey := fmt.Sprintf(db.OKTO_AUTH_TOKEN_KEY, chatId)
	authTokenStr, err := db.RedisClient().Get(context.Background(), authTokenKey).Result()
	if err != nil {
		return "", errors.Wrap(err, "error fetching okto auth token")
	}
	// TODO: handle token not found
	authToken := okto.AuthToken{}
	// TODO: handle decoding error
	json.NewDecoder(strings.NewReader(authTokenStr)).Decode(&authToken)
	return authToken.AuthToken, nil
}

func GetAddress(chatId int64, network string) (string, error) {
	addressKey := fmt.Sprintf(db.OKTO_ADDRESSES_KEY, chatId)
	addrRes, err := db.RedisClient().Get(context.Background(), addressKey).Result()
	if err != nil {
		return "", errors.Wrap(err, "error getting addresses from redis")
	}

	var wallets []okto.Wallet
	err = json.NewDecoder(strings.NewReader(addrRes)).Decode(&wallets)
	if err != nil {
		return "", errors.Wrap(err, "error decoding wallets")
	}

	for _, w := range wallets {
		if w.NetworkName == network {
			return w.Address, nil
		}
	}

	return "", errors.New("network not found")
}
