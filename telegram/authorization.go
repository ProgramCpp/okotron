package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/google"
	"github.com/programcpp/okotron/okto"
)

func Login(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	// TODO: validate if already authorized. auth token would be removed after expiry. login again
	// dbKey := fmt.Sprintf("%d", id)
	// token := db.Get(dbKey)
	// // TODO: do proper error check. check error from db call
	// if token != "" {
	// 	Send(bot, update, "okotron already authorized.")
	// 	return
	// }

	deviceCode, err := google.GetDeviceCode()
	if err != nil {
		log.Printf("error encountered when saving token id %d", id)
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}
	reply := ""
	reply += fmt.Sprintf("visit [google authorization page](%s) to authorize okotron and enter device code %s. \n", deviceCode.VerificationUrl, deviceCode.UserCode)
	reply += "return to okotron chat when done"
	Send(bot, update, reply)

	for i := 0; i*deviceCode.Interval <= deviceCode.ExpiresIn; i++ {
		time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
		googleToken, err := google.PollAuthorization(deviceCode.DeviceCode)
		if errors.Is(err, google.ErrAuthorizationPending) {
			continue
		} else if err != nil {
			log.Printf("error encountered when polling for token id for user %d", id)
			Send(bot, update, "error authorizing okotron. try again.")
			break
		} else {
			authenticate(bot, update, googleToken)
			break
		}
	}
}

func authenticate(bot *tgbotapi.BotAPI, update tgbotapi.Update, googleToken google.AccessToken) {
	id := update.Message.Chat.ID
	googleIDTokenKey := fmt.Sprintf(db.GOOGLE_ID_TOKEN_KEY, id) // TODO: expire tokens
	err := db.RedisClient().Set(context.Background(), googleIDTokenKey, googleToken.IdToken, 0).Err()
	if err != nil {
		log.Printf("error encountered when saving google token id %d. %s", id, err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	token, err := okto.Authenticate(googleToken.IdToken)
	if err != nil {
		log.Println("error authentication to Okto. " + err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	// TODO: store expiry for refresh
	tokenKey := fmt.Sprintf(db.OKTO_AUTH_TOKEN_KEY, id)
	err = db.RedisClient().Set(context.Background(), tokenKey, token, 0).Err()
	if err != nil {
		log.Printf("error encountered when saving token id %d. %s", id, err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	// TODO: create wallets only if not already created. enquire wallets
	wallets, err := okto.CreateWallet(token.AuthToken)
	if err != nil {
		log.Println("error authentication to Okto. " + err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(wallets)
	if err != nil {
		log.Println("error encoding Okto wallets. " + err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	addressKey := fmt.Sprintf(db.OKTO_ADDRESSES_KEY, update.Message.Chat.ID)
	err = db.RedisClient().Set(context.Background(), addressKey, buf.Bytes(), 0).Err()
	if err != nil {
		log.Println("error saving okto wallets. " + err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	reply := "okotron setup is now complete. fund your wallets to get started \n"

	// display wallets for users to fund them
	for _, w := range wallets {
		reply += fmt.Sprintf("Network: %s \n Wallet address: %s\n\n", w.NetworkName, w.Address)
	}

	Send(bot, update, reply)
}
