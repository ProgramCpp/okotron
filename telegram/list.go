package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/db"
	"github.com/programcpp/oktron/okto"
)

func List(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	dbKey := fmt.Sprintf("%d", id)
	authToken := db.Get(dbKey)
	okto.SupportedTokens(authToken)
}
