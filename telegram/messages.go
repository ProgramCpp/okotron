package telegram

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/oktron/db"
)

type COMMANDS string

// primary commands
const (
	CMD_LOGIN     = "/login"
	CMD_PORTFOLIO = "/portfolio"
	CMD_SWAP      = "/swap"
)

// sub commands that can be executed after the primary commands
const (
	// TODO: fix the weird naming convention to connect commands and sub commands
	CMD_LOGIN_CMD_SETUP_PROFILE = "/login/setup-profile"
	CMD_SWAP_CMD_SELECT_SOURCE  = "/swap/select-source"
)

// blocking call. reads telegram messages and processes them
func Run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic("telegram bot token missing:", err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	debugLevel, err := strconv.ParseBool(os.Getenv("ENABLE_DEBUG_LOGS"))
	if err != nil {
		log.Println("invalid value for config ENABLE_DEBUG_LOGS")
	}

	bot.Debug = debugLevel

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// fmt.Printf("%+v", update)
		if update.Message != nil { // If we got a message
			command := update.Message.Text
			subCommand := ""
			if update.Message.ReplyToMessage != nil {
				messageKey := fmt.Sprintf("message_%d", update.Message.MessageID)
				subCommand = db.Get(messageKey)
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if command == CMD_LOGIN {
				go Login(bot, update)
			} else if command == CMD_PORTFOLIO {
				go Portfolio(bot, update)
			} else if command == CMD_SWAP {
				go Swap(bot, update)
			} else if subCommand == CMD_LOGIN_CMD_SETUP_PROFILE {
				go SetupProfile(bot, update)
			} else {
				go Greet(bot, update)
			}
		} else if update.CallbackQuery != nil {
			// todo: command to go back
			if update.CallbackQuery.Data == CMD_SWAP_CMD_SELECT_SOURCE {
				go SwapSourceNetwork(bot, update)
			}
		}
	}
}
