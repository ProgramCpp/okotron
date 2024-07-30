package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/db"
	"github.com/programcpp/oktron/okto"
)

func SetupProfile(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	tokenKey := fmt.Sprintf("okto_token_%d", id)
	token := db.Get(tokenKey)
	googleIdTokenKey := fmt.Sprintf("google_id_token_%d", id)
	idToken := db.Get(googleIdTokenKey)
	authToken, err := okto.SetPin(idToken, token, update.Message.Text)
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

	authTokenKey := fmt.Sprintf("okto_auth_token_%d", id)
	db.Save(authTokenKey, buffer.String())

	wallets, err := okto.CreateWallet(authToken.AuthToken)
	if err != nil {
		log.Println("error authentication to Okto. " + err.Error())
		Send(bot, update, "error authorizing oktron. try again.")
		return
	}
	reply := "oktron setup is now complete. fund your wallets to get started \n"

	// display wallets for users to fund them
	for _, w := range wallets {
		reply += fmt.Sprintf("network: %s \n wallet address: %s\n\n", w.NetworkName, w.Address)
	}

	Send(bot, update, reply)
}
