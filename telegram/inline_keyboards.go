package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	tokenKeyboard = func() tgbotapi.InlineKeyboardMarkup {
		keyboardButtons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("back", "back"),
		}

		for _, n := range SUPPORTED_TOKENS {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
		}

		var tokenKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
		)

		return tokenKeyboard
	}()

	// this is just a function - a light weight function. only initializes the keyboard
	// or just a variable dynamically initialized
	// note the slight difference to the above token keyboard. tokenKeyboard is a keyboard. networkKeyboard is a func that creates a keyboard.
	// all other keyboards are created once except network keyboard, which is created dynamically based on token for which networks must be listed
	// hope I was clear :P I know this is not consistent. but once you build your mental model around it, it should be straightforward - you have a keyboard or a function that creates a keyboard -simple, right?
	networkKeyboard = func(toToken string) tgbotapi.InlineKeyboardMarkup {
		keyboardButtons := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("back", "back"),
		}

		for _, n := range SUPPORTED_NETWORKS[toToken] {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(n, n))
		}

		var networkKeyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
		)

		return networkKeyboard
	}

	numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
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
			tgbotapi.NewInlineKeyboardButtonData("back", "back"),
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("enter", "enter"),
		),
	)
)
