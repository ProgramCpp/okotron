package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/db"
	"github.com/programcpp/oktron/okto"
)

func Portfolio(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	authTokenKey := fmt.Sprintf("okto_auth_token_%d", id)
	authTokenStr := db.Get(authTokenKey)

	authToken := okto.AuthToken{}
	json.NewDecoder(strings.NewReader(authTokenStr)).Decode(&authToken)
	tokens, err := okto.SupportedTokens(authToken.AuthToken)
	// TODO: handle authorization failures. send descriptive message for user
	if err != nil {
		log.Printf("error fetching supported tokens. " + err.Error())
		Send(bot, update, "something went wrong!")
		return
	}
	reply := ""
	for _, token := range tokens {
		reply += token.String() + "\n"
	}
	Send(bot, update, reply)
}
