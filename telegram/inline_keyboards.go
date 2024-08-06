package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func tokenKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("back", "back"),
	}

	for _, n := range SUPPORTED_TOKENS {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	tokenKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	return tokenKeyboard
}

func networkKeyboard(toToken string) tgbotapi.InlineKeyboardMarkup {
	keyboardButtons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("back", "back"),
	}

	for _, n := range SUPPORTED_NETWORKS[toToken] {
		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
	}

	networkKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
	)

	return networkKeyboard
}

func numericKeyboard(back bool) tgbotapi.InlineKeyboardMarkup {
	var lastRow []tgbotapi.InlineKeyboardButton

	if back {
		lastRow = append(lastRow, tgbotapi.NewInlineKeyboardButtonData("back", "back"))
	}
	lastRow = append(lastRow,
		tgbotapi.NewInlineKeyboardButtonData("0", "0"),
		tgbotapi.NewInlineKeyboardButtonData("enter", "enter"),
	)

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
		lastRow)

	return numericKeyboard
}
