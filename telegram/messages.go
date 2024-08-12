package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
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
				messageKey := fmt.Sprintf(db.MESSAGE_KEY, update.Message.MessageID)
				subCommand, err = db.RedisClient().Get(context.Background(), messageKey).Result()
				if err != nil {
					Send(bot, update, "something went wrong. try again")
					continue
				}
			}

			if command == CMD_LOGIN {
				go Login(bot, update)
			} else if command == CMD_PORTFOLIO {
				go Portfolio(bot, update)
			} else if command == CMD_SWAP {
				go TokenInput(bot, update)
			} else if subCommand == CMD_LOGIN_CMD_SETUP_PROFILE {
				go SetupProfile(bot, update)
			} else if command == CMD_LIMIT_ORDER {
				go LimitOrder(bot, update)
			} else {
				go Greet(bot, update)
			}
		} else if update.CallbackQuery != nil {
			// todo: command to go back
			messageKey := fmt.Sprintf(db.MESSAGE_KEY, update.CallbackQuery.Message.MessageID)
			subCommand, err := db.RedisClient().Get(context.Background(), messageKey).Result()
			if err != nil {
				Send(bot, update, "something went wrong. try again")
				continue
			}
			isBack := update.CallbackQuery.Data == "back"
			if subCommand == CMD_ANY_CMD_FROM_TOKEN {
				go FromToken(bot, update, isBack)
			} else if subCommand == CMD_ANY_CMD_FROM_NETWORK {
				go FromNetwork(bot, update, isBack)
			} else if subCommand == CMD_ANY_CMD_TO_TOKEN {
				go ToToken(bot, update, isBack)
			} else if subCommand == CMD_ANY_CMD_TO_NETWORK {
				go ToNetwork(bot, update, isBack)
			} else if subCommand == CMD_ANY_CMD_QUANTITY {
				go Quantiy(bot, update, isBack)
			}
		}
	}
}
