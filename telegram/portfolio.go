package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/programcpp/okotron/okto"
)

func Portfolio(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	authTokenKey := fmt.Sprintf(db.OKTO_AUTH_TOKEN_KEY, id)
	authTokenStr, err := db.RedisClient().Get(context.Background(), authTokenKey).Result()
	if err != nil {
		log.Printf("error fetching okto auth token. " + err.Error())
		Send(bot, update, "something went wrong!")
		return
	}
	// TODO: handle token not found

	authToken := okto.AuthToken{}
	json.NewDecoder(strings.NewReader(authTokenStr)).Decode(&authToken)
	tokens, err := okto.Portfolio(authToken.AuthToken)
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
