package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendWithForceReply(bot *tgbotapi.BotAPI, update tgbotapi.Update, reply string, forceReply bool) (tgbotapi.Message, error) {
	const PARSE_MODE = "MarkdownV2"
	//reply = tgbotapi.EscapeText(PARSE_MODE, reply) // doesn't work for links!
	specialCharacterEscaper := strings.NewReplacer(
		`_`, "\\_",
		`*`, "\\*",
		// `[`, "\\[", // doesn't work for links!
		// `]`, "\\]",
		// `(`, "\\(",
		// `)`, "\\)",
		`~`, "\\~",
		"`", "\\`",
		`>`, "\\>",
		`#`, "\\#",
		`+`, "\\+",
		`-`, "\\-",
		`=`, "\\=",
		`|`, "\\|",
		`{`, "\\{",
		`}`, "\\}",
		`.`, "\\.",
		`!`, "\\!",
	)
	reply = specialCharacterEscaper.Replace(reply)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
	msg.ReplyToMessageID = update.Message.MessageID
	// https://core.telegram.org/bots/api#formatting-options
	msg.ParseMode = PARSE_MODE
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply: forceReply,
	}
	resp, err := bot.Send(msg)
	if err != nil {
		log.Println("error sending message to bot: " + err.Error())
		return tgbotapi.Message{}, err
	}
	return resp, nil
}

func Send(bot *tgbotapi.BotAPI, update tgbotapi.Update, reply string) (tgbotapi.Message, error) {
	return SendWithForceReply(bot, update, reply, false)
}
