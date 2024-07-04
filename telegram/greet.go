package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Greet(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	reply := "hello " + "first_name" // TODO: add default options
	Send(bot, update, reply)
}
