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
	deviceCode := google.GetDeviceCode()
	reply := ""
	reply += fmt.Sprintf("visit [google authorization page](%s) to authorize oktron. enter device code %s \n", deviceCode.Verification_url, deviceCode.DeviceCode)
	reply += "return to oktron chat when done"
	Send(bot, update, reply)

	for {
		time.Sleep(time.Duration(deviceCode.Interval) * time.Second)
		token, err := google.PollAuthorization()
		if errors.Is(err, google.ErrAuthorizationPending) {
			continue
		} else if err != nil {
			log.Printf("error encountered when polling for token id for user %d", id)
			Send(bot, update, "error authorizing oktron. try again.")
			break
		} else {
			err = db.Save(fmt.Sprintf("%d", id), token)
			if err != nil {
				log.Printf("error encountered when saving token id %d", id)
				Send(bot, update, "error authorizing oktron. try again.")
				break
			}

			Send(bot, update, "Authorization success")
		}
	}
}
