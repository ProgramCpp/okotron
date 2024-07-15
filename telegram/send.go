package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Send(bot *tgbotapi.BotAPI, update tgbotapi.Update, reply string){
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	// https://core.telegram.org/bots/api#formatting-options
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}