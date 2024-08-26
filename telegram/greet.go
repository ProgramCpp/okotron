package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Greet(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	reply := "hello " + update.Message.From.FirstName + ". please enter a valid command"
	 // TODO: add default options
	 // by default /start command must be handled. list commands or use an inline keyboard
	Send(bot, update, reply)
}
