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

const (
	LOGIN         = "/login"
	SETUP_PROFILE = "/setup-profile"
	PORTFOLIO     = "/portfolio"
	SWAP          = "/swap"
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
		command := update.Message.Text
		subCommand := ""
		if update.Message.ReplyToMessage != nil {
			messageKey := fmt.Sprintf("message_%d", update.Message.MessageID)
			subCommand = db.Get(messageKey)
		}
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if command == LOGIN {
				go Login(bot, update)
			} else if command == PORTFOLIO {
				go Portfolio(bot, update)
			} else if command == SWAP {
				go Swap(bot, update)
			} else if subCommand == SETUP_PROFILE {
				go SetupProfile(bot, update)
			} else {
				go Greet(bot, update)
			}
		}
	}
}
