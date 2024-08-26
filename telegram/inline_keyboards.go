package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/utils"
)

func BuyOrSellKeyboard() tgbotapi.InlineKeyboardMarkup {
	tokenKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("buy", "buy"),
			tgbotapi.NewInlineKeyboardButtonData("sell", "sell"),
		),
	)

	return tokenKeyboard
}

func tokenKeyboard(back bool) tgbotapi.InlineKeyboardMarkup {
	noOfButtonsPerRow := 2
	keyboardRows := [][]tgbotapi.InlineKeyboardButton{}

	if back {
		keyboardRows = append(keyboardRows, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("⬅back", "back")})
	}
	for i := 0; i < len(utils.SUPPORTED_TOKENS); {
		keyboardButtons := []tgbotapi.InlineKeyboardButton{}
		for j := 0; j < noOfButtonsPerRow && i < len(utils.SUPPORTED_TOKENS); j++ {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(utils.SUPPORTED_TOKENS[i], utils.SUPPORTED_TOKENS[i]))
			i++
		}
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(keyboardButtons...))
	}

	tokenKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		keyboardRows...,
	)

	return tokenKeyboard
}

func networkKeyboard(toToken string) tgbotapi.InlineKeyboardMarkup {
	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅back", "back"),
	}

	for _, n := range utils.SUPPORTED_NETWORKS[toToken] {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	networkKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	return networkKeyboard
}

func numericKeyboard() tgbotapi.InlineKeyboardMarkup {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("7", "7"),
			tgbotapi.NewInlineKeyboardButtonData("8", "8"),
			tgbotapi.NewInlineKeyboardButtonData("9", "9"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(".", "."),
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("enter ↩", "enter"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅back", "back"),
		),
	)

	return numericKeyboard
}
