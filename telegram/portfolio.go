package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/okto"
	"github.com/programcpp/okotron/utils"
)

func Portfolio(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID

	authToken, err := utils.GetAuthToken(id)
	if err != nil {
		log.Printf("error fetching okto auth token. " + err.Error())
		Send(bot, update, "something went wrong!")
		return
	}

	tokens, err := okto.Portfolio(authToken)
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
