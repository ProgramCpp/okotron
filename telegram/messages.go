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

type COMMANDS string

// primary commands
const (
	CMD_LOGIN       = "/login"
	CMD_PORTFOLIO   = "/portfolio"
	CMD_SWAP        = "/swap"
	CMD_LIMIT_ORDER = "/limit-order"
)

// sub commands that can be executed after the primary commands
const (
	// TODO: fix the weird naming convention to pass context from commands and sub commands
	CMD_LOGIN_CMD_SETUP_PROFILE = "/login/setup-profile"
	CMD_SWAP_CMD_FROM_TOKEN     = "/swap/from-token"
	CMD_SWAP_CMD_FROM_NETWORK   = "/swap/from-network"
	CMD_SWAP_CMD_TO_TOKEN       = "/swap/to-token"
	CMD_SWAP_CMD_TO_NETWORK     = "/swap/to-network"
	CMD_SWAP_CMD_QUANTITY       = "/swap/quantity"
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
				subCommand = db.Get(messageKey)
			}

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
			messageKey := fmt.Sprintf(db.MESSAGE_KEY, update.CallbackQuery.Message.MessageID)
			subCommand, err := db.RedisClient().Get(context.Background(), messageKey).Result()
			if err != nil {
				Send(bot, update, "something went wrong. try again")
				continue
			}
			isBack := update.CallbackQuery.Data == "back"
			if subCommand == CMD_SWAP_CMD_FROM_TOKEN {
				go SwapFromToken(bot, update, isBack)
			} else if subCommand == CMD_SWAP_CMD_FROM_NETWORK {
				go SwapFromNetwork(bot, update, isBack)
			} else if subCommand == CMD_SWAP_CMD_TO_TOKEN {
				go SwapToToken(bot, update, isBack)
			} else if subCommand == CMD_SWAP_CMD_TO_NETWORK {
				go SwapToNetwork(bot, update, isBack)
			} else if subCommand == CMD_SWAP_CMD_QUANTITY {
				go SwapQuantiy(bot, update, isBack)
			}
		}
	}
}
