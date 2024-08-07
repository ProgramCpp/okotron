package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
)

// sets pin and creates wallets
// call this command to finish authorization
func SetupProfile(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	pin := update.Message.Text
	tokenKey := fmt.Sprintf(db.OKTO_TOKEN_KEY, id)
	token := db.Get(tokenKey)
	googleIdTokenKey := fmt.Sprintf(db.GOOGLE_ID_TOKEN_KEY, id)
	idToken := db.Get(googleIdTokenKey)
	authToken, err := okto.SetPin(idToken, token, pin)
	if err != nil {
		log.Println("error setting okto pin" + err.Error())
		Send(bot, update, "encountered a problem when setting the PIN. try again.")
		return
	}

	buffer := bytes.Buffer{}
	err = json.NewEncoder(&buffer).Encode(authToken)
	if err != nil {
		log.Println("error serializing auth token" + err.Error())
		Send(bot, update, "encountered a problem when setting the PIN. try again.")
		return
	}

	authTokenKey := fmt.Sprintf(db.OKTO_AUTH_TOKEN_KEY, id)
	err = db.RedisClient().Set(context.Background(), authTokenKey, buffer.String(), 0).Err()
	if err != nil {
		log.Println("error saving auth token" + err.Error())
		Send(bot, update, "encountered a problem when setting the PIN. try again.")
		return
	}

	// TODO: create wallets only if not already created. enquire wallets
	wallets, err := okto.CreateWallet(authToken.AuthToken)
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

	err = db.RedisClient().Set(context.Background(), fmt.Sprintf(db.OKTO_ADDRESSES_KEY, update.Message.Chat.ID), buf.String(), 0).Err()
	if err != nil {
		log.Println("error saving okto walletso. " + err.Error())
		Send(bot, update, "error authorizing okotron. try again.")
		return
	}

	reply := "okotron setup is now complete. fund your wallets to get started \n"

	// display wallets for users to fund them
	for _, w := range wallets {
		reply += fmt.Sprintf("network: %s \n wallet address: %s\n\n", w.NetworkName, w.Address)
	}

	Send(bot, update, reply)
}
