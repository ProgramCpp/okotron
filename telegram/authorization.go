package telegram

import (
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/programcpp/oktron/db"
	"github.com/programcpp/oktron/google"
)

func Auth(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	id := update.Message.Chat.ID
	deviceCode, err := google.GetDeviceCode()
	if err != nil {
		log.Printf("error encountered when saving token id %d", id)
		Send(bot, update, "error authorizing oktron. try again.")
		return
	}
	reply := ""
	reply += fmt.Sprintf("visit [google authorization page](%s) to authorize oktron. enter device code %s \n", deviceCode.VerificationUrl, deviceCode.UserCode)
	reply += "return to oktron chat when done"
	Send(bot, update, reply)

	for i := 0; i*deviceCode.Interval <= deviceCode.ExpiresIn; i++ {
		time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
		token, err := google.PollAuthorization(deviceCode.DeviceCode)
		if errors.Is(err, google.ErrAuthorizationPending) {
			continue
		} else if err != nil {
			log.Printf("error encountered when polling for token id for user %d", id)
			Send(bot, update, "error authorizing oktron. try again.")
			break
		} else {
			err = db.Save(fmt.Sprintf("%d", id), token.AccessToken)
			if err != nil {
				log.Printf("error encountered when saving token id %d", id)
				Send(bot, update, "error authorizing oktron. try again.")
				break
			}

			Send(bot, update, "Authorization success")
		}
	}
}
