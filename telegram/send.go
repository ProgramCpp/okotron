package telegram

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Send(bot *tgbotapi.BotAPI, update tgbotapi.Update, reply string) error {
	specialCharacterEscaper := strings.NewReplacer(
		`_`, "\\_",
		`*`, "\\*",
		// `[`, "\\[",
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
	msg.ParseMode = "MarkdownV2"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("error sending message to bot: " + err.Error())
		return err
	}
	return nil
}
