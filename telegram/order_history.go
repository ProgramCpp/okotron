package telegram

import (
	"errors"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/okto"
)

func OrderHistory(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID

	authToken, err := okto.GetAuthToken(id)
	if err != nil {
		log.Printf("error fetching okto auth token. " + err.Error())
		Send(bot, update, "something went wrong!")
		return
	}

	jobs, err := okto.OrderHistory(authToken)
	// TODO: handle authorization failures. send descriptive message for user
	if err != nil && errors.Is(err, okto.ERR_UNAUTHORIZED){
		log.Printf("error fetching order history. " + err.Error())
		Send(bot, update, "unauthorized. login and try again.")
		return
	} else if err != nil {
		log.Printf("error fetching order history. " + err.Error())
		Send(bot, update, "something went wrong!")
		return
	}
	reply := ""
	for _, job := range jobs {
		reply += job.String() + "\n"
	}
	Send(bot, update, reply)
}