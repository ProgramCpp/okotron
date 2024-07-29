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
	LOGIN  = "/login"
	PIN    = "/pin"
	TOKENS = "/tokens"
	BUY    = "/buy"
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
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.Text == LOGIN {
				go Login(bot, update)
			} else if update.Message.Text == TOKENS {
				go Tokens(bot, update)
			} else if update.Message.Text == BUY {
				go Buy(bot, update)
			} else {
				// handle sub commands
				if update.Message.ReplyToMessage != nil {
					messageKey := fmt.Sprintf("message_%d", update.Message.MessageID)
					command := db.Get(messageKey)
					if command == PIN {
						go SetupProfile(bot, update)
						continue
					}
				}
				go Greet(bot, update)
			}
		}
	}
}
