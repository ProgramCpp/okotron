package telegram

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/programcpp/okotron/db"
	"github.com/spf13/viper"
)

func CopyTradeInput(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// TODO: handle error
	resp, _ := bot.Send(tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: update.FromChat().ID,
			// TODO: consolidate all telegram send messages
			ReplyMarkup: tgbotapi.ForceReply{
				ForceReply: true,
			},
		},
		Text: "enter the address:",
	})
	subcommandKey := fmt.Sprintf(db.SUB_COMMAND_KEY, resp.MessageID)
	err := db.RedisClient().Set(context.Background(), subcommandKey, CMD_COPY_TRADE_CMD_ADDRESS,
		time.Duration(viper.GetInt("REDIS_CMD_EXPIRY_IN_SEC"))*time.Second).Err()
	if err != nil {
		log.Printf("error encountered when saving sub command key for copy trade. %s", err.Error())
		bot.Send(tgbotapi.NewEditMessageText(update.FromChat().ID, resp.MessageID, "something went wrong. try again."))
	}
}

func CopyTradeAddressInput(bot *tgbotapi.BotAPI, update tgbotapi.Update, isBack bool) {
	address := update.Message.Text

	copyOrderKey := fmt.Sprintf(db.COPY_ORDER_KEY, address)
	// TODO: handle error
	db.RedisClient().RPush(context.Background(), copyOrderKey, fmt.Sprintf("%d", update.FromChat().ID))

	bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "copy trade success"))
}

func CopyTrade(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	CopyTradeInput(bot, update)
}

func CopyTradeCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

}