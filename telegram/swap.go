package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// TODO: use okto /aupported_tokens api
	// for now all networks returnd by supported networks do not work. ex: solana
	// an array. do not handle each network separately. oktron is network agnostic
	SUPPORTED_NETWORKS = []string{"APTOS", "BASE", "POLYGON"}
)

func Swap(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: separate UI componets from backend
	keyboardButtons := []tgbotapi.InlineKeyboardButton{}

	for _, n := range SUPPORTED_NETWORKS {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID : update.Message.Chat.ID,
			ReplyMarkup: networkKeyboard,
		},
		Text: "select the source network",
	})
}
